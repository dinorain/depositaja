package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"

	"github.com/dinorain/depositaja"
	"github.com/dinorain/depositaja/collector"
	"github.com/dinorain/depositaja/flagger"
	"github.com/dinorain/depositaja/proto/pb"
)

type depositRequest struct {
	WalletID string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
}

type checkResponse struct {
	WalletID       string  `json:"wallet_id"`
	Balance        float64 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Run(brokers []string, stream goka.Stream) {
	view, err := goka.NewView(brokers, collector.Table, new(depositaja.DepositListCodec))
	if err != nil {
		panic(err)
	}
	go view.Run(context.Background())

	flaggerView, err := goka.NewView(brokers, flagger.Table, new(flagger.FlagValueCodec))
	if err != nil {
		panic(err)
	}
	go flaggerView.Run(context.Background())

	emitter, err := goka.NewEmitter(brokers, stream, new(depositaja.DepositCodec))
	if err != nil {
		panic(err)
	}
	defer emitter.Finish()

	router := mux.NewRouter()
	router.HandleFunc("/deposit", deposit(emitter, stream)).Methods("POST")
	router.HandleFunc("/check/{wallet_id}", check(view, flaggerView)).Methods("GET")

	log.Printf("Listen port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func deposit(emitter *goka.Emitter, stream goka.Stream) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req depositRequest

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			//log.Printf("error io: %v", err)
			respondWithError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		err = json.Unmarshal(b, &req)
		if err != nil {
			//log.Printf("error json unmarshal: %v", err)
			respondWithError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if !(req.Amount > 0) {
			//log.Printf("error invalid amount: %v", deposit.Amount)
			respondWithError(w, http.StatusUnprocessableEntity, "amount must be more than 0")
			return
		}

		deposit := &pb.Deposit{
			WalletId: req.WalletID,
			Amount:   req.Amount,
		}

		if stream == depositaja.DepositStream {
			err = emitter.EmitSync(req.WalletID, deposit)
		} else {
			deposit.Amount = -1 * deposit.Amount
			err = emitter.EmitSync(req.WalletID, deposit)
		}
		if err != nil {
			//log.Printf("error emit: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		//log.Printf("Wallet statement:\n %v to %v\n", deposit.Amount, deposit.WalletID)
		respondWithJSON(w, http.StatusOK, nil)
	}
}

func check(view *goka.View, flaggerView *goka.View) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walletID := mux.Vars(r)["wallet_id"]

		var totalBalance float64
		var aboveThreshold bool

		response := checkResponse{
			WalletID:       walletID,
			Balance:        totalBalance,
			AboveThreshold: aboveThreshold,
		}

		val, _ := view.Get(walletID)
		if val == nil {
			//log.Printf("%s not found!", walletID)
			respondWithJSON(w, http.StatusOK, response)
			return
		}

		// log.Printf("Wallet statement for %s\n", walletID)
		messages := val.(*pb.DepositHistory)
		for _, m := range messages.Deposits {
			totalBalance += m.Amount
			// log.Printf("%d %10s: %v\n", i, m.WalletID, m.Amount)
		}

		flaggerVal, _ := flaggerView.Get(walletID)
		if flaggerVal != nil {
			b := flaggerVal.(*pb.FlagValue)
			aboveThreshold = b.Flagged
		}

		// log.Printf("Balance: %v\n", totalBalance)
		response.Balance = totalBalance
		response.AboveThreshold = aboveThreshold
		respondWithJSON(w, http.StatusOK, response)
	}
}
