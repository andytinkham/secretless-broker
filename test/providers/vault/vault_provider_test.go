package main

import (
	"testing"

	_ "github.com/joho/godotenv/autoload"
	. "github.com/smartystreets/goconvey/convey"

	plugin_v1 "github.com/cyberark/secretless-broker/internal/plugin/v1"
	"github.com/cyberark/secretless-broker/internal/providers"
)

func TestVault_Provider(t *testing.T) {
	var err error
	var provider plugin_v1.Provider
	name := "vault"

	options := plugin_v1.ProviderOptions{
		Name: name,
	}

	Convey("Can create the Vault provider", t, func() {
		provider, err = providers.ProviderFactories[name](options)
		So(err, ShouldBeNil)
	})

	Convey("Has the expected provider name", t, func() {
		So(provider.GetName(), ShouldEqual, "vault")
	})

	Convey("Reports when the secret is not found", t, func() {
		id := "foobar"

		values, err := provider.GetValues(id)
		So(err, ShouldBeNil)
		So(values[id], ShouldNotBeNil)
		So(values[id].Error, ShouldNotBeNil)
		So(values[id].Error.Error(), ShouldEqual, "HashiCorp Vault provider could not find a secret called 'foobar'")
		So(values[id].Value, ShouldBeNil)
	})

	Convey("Can provide a secret", t, func() {
		id := "kv/db/password#password"
		values, err := provider.GetValues(id)

		So(err, ShouldBeNil)
		So(values[id], ShouldNotBeNil)
		So(values[id].Error, ShouldBeNil)
		So(values[id].Value, ShouldNotBeNil)
		So(string(values[id].Value), ShouldEqual, "db-secret")
	})

	Convey("Can provide a secret with default field name", t, func() {
		id := "kv/web/password"
		values, err := provider.GetValues(id)

		So(err, ShouldBeNil)
		So(values[id], ShouldNotBeNil)
		So(values[id].Error, ShouldBeNil)
		So(values[id].Value, ShouldNotBeNil)
		So(string(values[id].Value), ShouldEqual, "web-secret")
	})
}
