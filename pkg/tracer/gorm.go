package tracer

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	//DBTraceContextKey opentracing在db中保存context的Key
	DBTraceContextKey = "trace.db.context"
	//DBTraceSpanKey opentracing在db中保存span的Key
	DBTraceSpanKey = "trace.db.span"
)

//GormTraceHandler opentracing处理gorm接口
type GormTraceHandler interface {
	CreateBefore(scope *gorm.Scope)
	CreateAfter(scope *gorm.Scope)
	QueryBefore(scope *gorm.Scope)
	QueryAfter(scope *gorm.Scope)
	RowQueryBefore(scope *gorm.Scope)
	RowQueryAfter(scope *gorm.Scope)
	UpdateBefore(scope *gorm.Scope)
	UpdateAfter(scope *gorm.Scope)
	DeleteBefore(scope *gorm.Scope)
	DeleteAfter(scope *gorm.Scope)
	RegisterDB(db *gorm.DB)
}

//gormTraceHandler opentracing默认处理gorm方法
type gormTraceHandler struct {
}

// RegisterDB .
func RegisterDB(db *gorm.DB) {
	//目前只支持全局配置 要么全部跟踪 要么全部不跟踪
	if !Enable {
		return
	}
	gth := gormTraceHandler{}
	//注册gorm操作数据库前后 trace事件
	db.Callback().Create().Before("gorm:create").Register("trace_plugin:before_create", gth.CreateBefore)
	db.Callback().Create().After("gorm:create").Register("trace_plugin:after_create", gth.CreateAfter)

	db.Callback().Query().Before("gorm:query").Register("trace_plugin:before_query", gth.QueryBefore)
	db.Callback().Query().After("gorm:query").Register("trace_plugin:after_query", gth.QueryAfter)

	db.Callback().Update().Before("gorm:update").Register("trace_plugin:before_update", gth.UpdateBefore)
	db.Callback().Update().After("gorm:update").Register("trace_plugin:after_update", gth.UpdateAfter)

	db.Callback().Delete().Before("gorm:delete").Register("trace_plugin:before_delete", gth.DeleteBefore)
	db.Callback().Delete().After("gorm:delete").Register("trace_plugin:after_delete", gth.DeleteAfter)

	db.Callback().RowQuery().Before("gorm:row_query").Register("trace_plugin:before_row_query", gth.RowQueryBefore)
	db.Callback().RowQuery().After("gorm:row_query").Register("trace_plugin:after_row_query", gth.RowQueryAfter)
}

//CreateBefore gorm before 处理
func (gth gormTraceHandler) CreateBefore(scope *gorm.Scope) {
	gth.Before(scope, "create")
}

//CreateAfter gorm after 处理
func (gth gormTraceHandler) CreateAfter(scope *gorm.Scope) {
	gth.After(scope, "create")
}

//QueryBefore gorm before 处理
func (gth gormTraceHandler) QueryBefore(scope *gorm.Scope) {
	gth.Before(scope, "query")
}

//QueryAfter gorm after 处理
func (gth gormTraceHandler) QueryAfter(scope *gorm.Scope) {
	gth.After(scope, "query")
}

//RowQueryBefore gorm before 处理
func (gth gormTraceHandler) RowQueryBefore(scope *gorm.Scope) {
	gth.Before(scope, "row_query")
}

//RowQueryAfter gorm after 处理
func (gth gormTraceHandler) RowQueryAfter(scope *gorm.Scope) {
	gth.After(scope, "row_query")
}

//UpdateBefore gorm before 处理
func (gth gormTraceHandler) UpdateBefore(scope *gorm.Scope) {
	gth.Before(scope, "update")
}

//UpdateAfter gorm after 处理
func (gth gormTraceHandler) UpdateAfter(scope *gorm.Scope) {
	gth.After(scope, "update")
}

//DeleteBefore gorm before 处理
func (gth gormTraceHandler) DeleteBefore(scope *gorm.Scope) {
	gth.Before(scope, "delete")
}

//DeleteAfter gorm after 处理
func (gth gormTraceHandler) DeleteAfter(scope *gorm.Scope) {
	gth.After(scope, "delete")
}

//Before opentracing通用的在sql执行前处理
func (gth gormTraceHandler) Before(scope *gorm.Scope, operation string) {
	if Enable {
		dataBaseName := scope.Dialect().CurrentDatabase()
		var operationName = dataBaseName + "." + operation
		var ctx context.Context
		if iCtx, iExists := scope.Get(DBTraceContextKey); iExists {
			ctx, _ = iCtx.(context.Context)
		}
		span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
		//设置tracer tag 数据库名称
		ext.DBInstance.Set(span, dataBaseName)
		//设置tracer tag 数据库类型
		ext.DBType.Set(span, "mssql")
		//然后把span传递出去
		scope.Set(DBTraceSpanKey, span)
	}
}

// DBStatement opentracing ext.DBStatement
var DBStatement = string(ext.DBStatement)

//After opentracing通用的在sql执行后处理
func (gth gormTraceHandler) After(scope *gorm.Scope, operation string) {
	if Enable {
		iSpan, iExists := scope.Get(DBTraceSpanKey)
		if !iExists {
			return
		}
		span, ok := iSpan.(opentracing.Span)
		if !ok || span == nil {
			return
		}

		ext.DBStatement.Set(span, scope.SQL)
		span.LogFields(log.String(DBStatement, scope.SQL))
		span.LogFields(log.Object("db.vars", scope.SQLVars))

		span.Finish()
	}
}

// NewTracerDB 获取新的带有span的DB
func NewTracerDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if !Enable {
		return db
	}
	return db.Set(DBTraceContextKey, ctx)
}
