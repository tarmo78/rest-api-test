/*
 * GUARDTIME CONFIDENTIAL
 *
 * Copyright 2008-2020 Guardtime, Inc.
 * All Rights Reserved.
 *
 * All information contained herein is, and remains, the property
 * of Guardtime, Inc. and its suppliers, if any.
 * The intellectual and technical concepts contained herein are
 * proprietary to Guardtime, Inc. and its suppliers and may be
 * covered by U.S. and foreign patents and patents in process,
 * and/or are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Guardtime, Inc.
 * "Guardtime" and "KSI" are trademarks or registered trademarks of
 * Guardtime, Inc., and no license to trademarks is granted; Guardtime
 * reserves and retains all trademark rights.
 */

package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	Controller struct {
	}

	Error struct {
		Error         string `json:"error"`
		ExtendedError string `json:"extendedError,omitempty"`
	}
)

func NewController() (*Controller, error) {

	controller := &Controller{}

	return controller, nil
}

func (c *Controller) SetupRouter(router *mux.Router) {
	apiRouter := router.PathPrefix("/test-service").Subrouter()
	apiRouter.Use(loggerMiddleware)
	apiRouter.HandleFunc("/webhook", c.testRequest).Methods(http.MethodPost)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading request body: %v\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reader := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader

		fmt.Printf("Server: request from=%s to=%s:%s\nheaders: %v\nbody: %v\n", r.RemoteAddr, r.Method, r.RequestURI,
			r.Header, string(buf))

		next.ServeHTTP(w, r)
	})
}
