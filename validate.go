package Validator

import (

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gogf/gf/util/gconv"
	"reflect"
)


const (
	Required = "required"
	RequiredMsgKey="requiredMsg"
	RequiredCodeKey="requiredCode"
)

type errorResult struct {
	ErrMsg string
	Code int
}

type c_validator struct {
	errDefaultTagMsg  map[string]errorTagMsg//主动设置的tag提示

}
//触发限制标签 提示语 比如触发  required 就提示 缺少必要参数（defaultMsg）
type errorTagMsg struct {
	defaultMsg  string
	defaultCode int
}

func NewValidator()  *c_validator{
	return &c_validator{map[string]errorTagMsg{}}
}

func (v *c_validator)SetErrMsg(msgType string,msg string)*c_validator {
	s:=v.errDefaultTagMsg[msgType]
	s.defaultMsg=msg
	v.errDefaultTagMsg[msgType]=s
	return v

}
func (v *c_validator)SetErrCode(msgType string,code int) *c_validator {
	s:=v.errDefaultTagMsg[msgType]
	s.defaultCode=code
	v.errDefaultTagMsg[msgType]=s
	return v
}


func (v *c_validator)GinContext(ctx *gin.Context,obj interface{})*errorResult {

	typeOfObj:=reflect.TypeOf(obj).Elem()
	if err:= ctx.ShouldBind(obj); err != nil {
		if fieldError,ok:=err.(validator.ValidationErrors);ok{
			fieldName:=fieldError[0].StructField()
			err1:=v.getErrorResult(fieldError,typeOfObj,fieldName)
			if(err1.Code==4008){
				err1.ErrMsg=err.Error()
			}
			return err1
		}
	}
	return nil
}

/*func (v *c_validator)Struct(obj interface{})*errorResult  {
	validate:= validator.New()
	err:= validate.Struct(obj)
	typeOfObj:=reflect.TypeOf(obj).Elem()
	if fieldError,ok:=err.(validator.ValidationErrors);ok{
		fieldName:=fieldError[0].StructField()
		err1:=v.getErrorResult(fieldError,typeOfObj,fieldName)
		if(err1.Code==4008){
			err1.ErrMsg=err.Error()
		}
		return err1
	}
	return nil
}*/


/**
比如需要增新的tag required_with
就switch 多一个 required_with,加一行如下代码
	return v.returnErrResult(field,fieldName,Required,RequiredMsgKey,RequiredCodeKey)
并且 const增加相应的常量
 */
func (v *c_validator) getErrorResult(fieldError validator.ValidationErrors, typeOfObj reflect.Type, fieldName string) (*errorResult) {

	field,_:=typeOfObj.FieldByName(fieldName);
	switch fieldError[0].Tag() {
		case "required":
			return v.returnErrResult(field,fieldName,Required,RequiredMsgKey,RequiredCodeKey)
			break
		default:
			return &errorResult{"",4008}
	}
	return nil
}


func (v *c_validator) returnErrResult(field reflect.StructField ,fieldName string,tagName string,tagErrmsgKey string,tagErrCodeKey string)  *errorResult {

	if msg:=v.errDefaultTagMsg[tagName].defaultMsg;msg!=""{
		code:=v.errDefaultTagMsg[tagName].defaultCode
		return &errorResult{msg+";field is "+fieldName,code}
	}else{
		msg=field.Tag.Get(tagErrmsgKey)
		code:=gconv.Int(field.Tag.Get(tagErrCodeKey))
		return &errorResult{msg+";field is "+fieldName,code}
	}
}
