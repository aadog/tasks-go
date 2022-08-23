package board

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"text/template"
)

const (
	Tpl_Count    = "{{.Count}}"
	Tpl_Success  = "{{.Success}}"
	Tpl_Error    = "{{.Error}}"
	Tpl_Exec     = "{{.Exec}}"
	Tpl_Noexec   = "{{.Noexec}}"
	Tpl_Process  = "{{.Process}}"
	Tpl_Process1 = "{{.Process1}}"
)

func Tpl_Custom(k interface{}) string {
	switch k.(type) {
	case string:
		return fmt.Sprintf("{{.Custom \"%v\" }}", k)
	}
	return fmt.Sprintf("{{.Custom %v }}", k)
}

// 计数板
type CountBoard struct {
	count    int32 //总数
	success  int32 //成功
	error    int32 //失败
	exec     int32 //已执行
	noexec   int32 //未执行
	template atomic.Value
	custom   sync.Map
}

func (this *CountBoard) SetCustom(custom sync.Map) {
	this.custom = custom
}

func (this *CountBoard) SetTemplate(template atomic.Value) {
	this.template = template
}

func (this *CountBoard) SetNoexec(noexec int32) {
	this.noexec = noexec
}

func (this *CountBoard) SetExec(exec int32) {
	this.exec = exec
}

func (this *CountBoard) SetError(error int32) {
	this.error = error
}

func (this *CountBoard) SetSuccess(success int32) {
	this.success = success
}
func (this *CountBoard) SetToStringtpl(templatestr string, vals ...interface{}) error {
	tmp, err := template.New("test").Parse(fmt.Sprintf(templatestr, vals...))
	if err != nil {
		return err
	}
	this.template.Store(tmp)
	return nil
}
func (this *CountBoard) Process() float32 {
	a := this.Exec()
	b := this.Count()
	c := float32(a*100/b) / 100
	return c
}
func (this *CountBoard) Process1() float32 {
	a := this.Exec()
	b := this.Count()
	c := float32(a * 100 / b)
	return c
}

func (this *CountBoard) Reset() {
	this.custom = sync.Map{}
	atomic.StoreInt32(&this.count, 0)
	atomic.StoreInt32(&this.success, 0)
	atomic.StoreInt32(&this.error, 0)
	atomic.StoreInt32(&this.exec, 0)
	atomic.StoreInt32(&this.noexec, 0)
}

func (this *CountBoard) Custom(k interface{}) int {
	iv, ok := this.custom.Load(k)
	if ok == false {
		v := int32(0)
		this.custom.Store(k, &v)
		return int(v)
	} else {
		v := iv.(*int32)
		return int(*v)
	}
}
func (this *CountBoard) Count() int {
	return int(atomic.LoadInt32(&this.count))
}
func (this *CountBoard) Success() int {
	return int(atomic.LoadInt32(&this.success))
}
func (this *CountBoard) Error() int {
	return int(atomic.LoadInt32(&this.error))
}
func (this *CountBoard) Exec() int {
	return int(atomic.LoadInt32(&this.exec))
}
func (this *CountBoard) Noexec() int {
	return int(atomic.LoadInt32(&this.noexec))
}
func (this *CountBoard) SetCount(count int) {
	this.Reset()
	atomic.StoreInt32(&this.count, int32(count))
	atomic.StoreInt32(&this.noexec, int32(count))
	atomic.StoreInt32(&this.success, 0)
	atomic.StoreInt32(&this.error, 0)
	atomic.StoreInt32(&this.exec, 0)
}
func (this *CountBoard) AddCustom(k interface{}, n int) {
	iv, ok := this.custom.Load(k)
	if ok == false {
		v := int32(n)
		this.custom.Store(k, &v)
	} else {
		v := iv.(*int32)
		atomic.AddInt32(v, int32(n))
	}
}
func (this *CountBoard) AddSuccess(n int) {
	for i := 0; i < n; i++ {
		this.Done()
	}
	atomic.AddInt32(&this.success, int32(n))
}
func (this *CountBoard) AddError(n int) {
	for i := 0; i < n; i++ {
		this.Done()
	}
	atomic.AddInt32(&this.error, int32(n))
}
func (this *CountBoard) AddSuccess_nodone(n int) {
	atomic.AddInt32(&this.success, int32(n))
}
func (this *CountBoard) AddError_nodone(n int) {
	atomic.AddInt32(&this.error, int32(n))
}
func (this *CountBoard) Done() {
	atomic.AddInt32(&this.exec, 1)
	atomic.AddInt32(&this.noexec, -1)
}
func (this *CountBoard) IsCompile() bool {
	return this.Process() >= 1
}

func (this *CountBoard) ToString() string {
	itl := this.template.Load()
	if itl == nil {
		return fmt.Sprintf("总数:%d    操作成功:%d     操作失败:%d     已执行:%d     未执行:%d",
			this.Count(),
			this.Success(),
			this.Error(),
			this.Exec(),
			this.Noexec(),
		)
	}
	tl := itl.(*template.Template)
	buf := bytes.Buffer{}
	err := tl.Execute(&buf, this)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func (this *CountBoard) ToTplString(templatestr string, vals ...interface{}) string {
	tmp, err := template.New("test").Parse(fmt.Sprintf(templatestr, vals...))
	if err != nil {
		return err.Error()
	}
	tl := tmp
	buf := bytes.Buffer{}
	err = tl.Execute(&buf, this)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func New(count int) *CountBoard {
	bd := &CountBoard{}
	bd.custom = sync.Map{}
	bd.SetCount(count)
	return bd
}
