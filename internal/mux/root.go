package mux

import (
	"html/template"
	"net/http"

	"github.com/pkg/errors"
	"github.com/yansal/requestr/internal/broker"
	"github.com/yansal/requestr/internal/store"
)

func rootHandler(store *store.Store, publisher broker.Publisher, template *template.Template) http.HandlerFunc {
	return handleError(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		if r.Method == "POST" {
			queue := r.PostFormValue("queue")
			payload := r.PostFormValue("payload")
			loop := r.PostFormValue("loop")
			if queue != "" && payload != "" {
				var publishOpts []broker.PublishOptions
				if loop != "" {
					publishOpts = append(publishOpts, broker.Loop)
				}
				if err := publisher.Publish(ctx, queue, payload, publishOpts...); err != nil {
					return err
				}
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return nil
		}

		pools, err := store.ListJobs(ctx)
		if err != nil {
			return err
		}
		return errors.WithStack(template.ExecuteTemplate(w, "index.html", pools))
	})
}
