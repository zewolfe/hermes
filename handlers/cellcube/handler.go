package cellcube

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/zewolfe/hermes/internal/log"
	"github.com/zewolfe/hermes/pkg/orchestrator"
	"github.com/zewolfe/hermes/pkg/rapidpro"
	"github.com/zewolfe/hermes/pkg/services/cellcube"
)

type handler struct {
	rp  *rapidpro.Service
	o   *orchestrator.Orchestrator
	log log.Logger
}

func Initialise(rp *rapidpro.Service) func(chi.Router) {
	o := orchestrator.New()
	h := handler{
		rp: rp,
		o:  o,
	}

	return func(r chi.Router) {
		r.Get("/in", h.handleIncomingFromMNO)
		r.Get("/in/*", h.handleIncomingFromMNO)
		r.Post("/send", h.handleIncomingFromRapidPro)
	}
}

func (h *handler) handleIncomingFromMNO(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	msisdn := queryParams.Get("MSISDN")
	path := r.URL.Path

	err := h.rp.TriggerFlow(msisdn, path)
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	ch := h.o.Subscribe(msisdn)
	data := <-ch

	switch data.(type) {
	case rapidpro.RapidProResponse:

		menu := &cellcube.Menu{}
		msg := data.(rapidpro.RapidProResponse)

		if err := json.Unmarshal([]byte(msg.Text), menu); err != nil {
			h.log.Error(err, "failed to unmarshal json", msg)
			http.Error(w, http.StatusText(500), 500)

			return
		}

		menuXml, err := menu.RenderXML()
		if err != nil {
			h.log.Error(err, "failed to render menu xml", menu)
			http.Error(w, http.StatusText(500), 500)

			return
		}

		go h.rp.Ack(msg, true)

		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(menuXml))

	default:
		//TODO: Handle
	}

}

// TODO: Find a better function name
func (h *handler) handleIncomingFromRapidPro(w http.ResponseWriter, r *http.Request) {
	data := &rapidpro.RapidProResponse{}
	if err := render.DecodeJSON(r.Body, data); err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	h.log.V(4).Info("Incoming data", "data", *data)

	h.o.Publish(data.To, *data)
	h.log.Info("Published - ", "ID", data.ID, "from", data.From, "to", data.To)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}
