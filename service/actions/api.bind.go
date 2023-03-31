package actions

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
	"service/service/info"
	"service/service/models"
	"strings"
)

func init() {
	_base.Subscribe(_base.TargetAction("api", "bind"), apiBind)
}

func apiBind(r uniform.IRequest, p diary.IPage) {
	var model models.Bind
	r.Read(&model)

	p.Notice("api.bind", diary.M{
		"model": model,
	})

	if err := p.Scope("http.bind", func(s diary.IPage) {
		s.Info("data", diary.M{
			"method": model.Method,
			"path":   model.Path,
		})
		info.Engine.Handle(model.Method, model.Path, _base.BindHandler(
			s,
			model.Timeout,
			model.Topic,
			nil,
			nil,
			nil,
		))

	}); err != nil {
		if !strings.HasPrefix(err.Error(), "handlers are already registered for path") {
			panic(err)
		}
	}

	if r.CanReply() {
		if err := r.Reply(uniform.Request{}); err != nil {
			p.Error("reply", err.Error(), diary.M{
				"err": err,
			})
		}
	}
}
