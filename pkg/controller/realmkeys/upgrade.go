// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package realmkeys

import (
	"net/http"

	"github.com/google/exposure-notifications-verification-server/pkg/controller"
)

func (c *Controller) HandleUpgrade() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		session := controller.SessionFromContext(ctx)
		if session == nil {
			controller.MissingSession(w, r, c.h)
			return
		}
		flash := controller.Flash(session)

		realm := controller.RealmFromContext(ctx)
		if realm == nil {
			controller.MissingRealm(w, r, c.h)
			return
		}

		if realm.CanUpgradeToRealmSigningKeys() {
			realm.UseRealmCertificateKey = true
			if err := c.db.SaveRealm(realm); err != nil {
				flash.Error("Error upgrading realm: %v", err)
				c.renderShow(ctx, w, r, realm)
				return
			} else {
				flash.Alert("Successfully switched to realm specific signing keys.")
			}
		} else {
			flash.Error("Issuer and Audience settings not complete, cannot upgrade.")
		}

		c.redirectShow(ctx, w, r)
	})
}
