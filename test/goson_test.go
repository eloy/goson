package goson_test

import (
	"github.com/harlock/goson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

type Address struct {
	City string
	State string
}

type User struct {
	Name string
	Id int
	Supervisor *User
	Addresses []Address
}

func (this User) Over() bool {
	return this.Id > 100
}

func (this User) Upper() string {
	return strings.ToUpper(this.Name)
}

var _ = Describe("Goson", func() {

	Describe("Creating a JSON object", func() {
		Context("Hash", func() {
			It("Should generage a json object", func() {
				foo := User{Name: "Foo", Id: 11 }
				g := goson.Hash(foo, "Id", "Over()")
				g.Method("Name")
				g.Alias("Uppercase()", "Upper()")
				data, err := g.ToJson()
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`{"id":11,"name":"Foo","over":false,"uppercase":"FOO"}`))
			})
		})

		Context("Array", func() {
			It("Should generage a json array", func() {
				foo := User{Name: "Foo", Id: 1}
				bar := User{Name: "Bar", Id: 200}
				array := []User{ foo, bar }

				g := goson.Array(array, "Id", "Name", "Over()", "Upper()")

				data, err := g.ToJson()
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`[{"id":1,"name":"Foo","over":false,"upper":"FOO"},{"id":200,"name":"Bar","over":true,"upper":"BAR"}]`))
			})
		})
	})

	// Nested Hash
	//----------------------------------------------------------------------

	Describe("Nested Hash", func() {
		Context("Hash", func() {
			It("Should include the nested hash", func() {
				foo := User{Name: "Foo", Id: 1}
				bar := User{Name: "Bar", Id: 200}
				wadus := User{Name: "Wadus", Id: 300}
				foo.Supervisor = &bar
				bar.Supervisor = &wadus

				g := goson.Hash(foo, "Id")
				supervisor := g.Hash("Supervisor", "Id").Method("Name")
				supervisor.HashAlias("Supervisor", "manager").Alias("manager_id", "Id")

				data, err := g.ToJson()
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`{"id":1,"supervisor":{"id":200,"manager":{"manager_id":300},"name":"Bar"}}`))
			})
		})

		Context("Array", func() {
			It("Should include the nested hash", func() {
				foo := User{Name: "Foo", Id: 1}
				bar := User{Name: "Bar", Id: 2}
				foowadus := User{Name: "FooWadus", Id: 301}
				barwadus := User{Name: "BarWadus", Id: 302}
				foo.Supervisor = &foowadus
				bar.Supervisor = &barwadus

				g := goson.Array([]User{foo, bar}, "Id")
				g.Hash("Supervisor", "Id").Method("Name")

				data, err := g.ToJson()
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`[{"id":1,"supervisor":{"id":301,"name":"FooWadus"}},{"id":2,"supervisor":{"id":302,"name":"BarWadus"}}]`))
			})
		})
	})



	// Nested Array
	//----------------------------------------------------------------------

	Describe("Nested Array", func() {
		Context("Hash", func() {
			It("Should include the nested array", func() {
				foo := User{Name: "Foo", Id: 1}
				address1 := Address{City:"Ctr", State:"al"}
				address2 := Address{City:"Tr", State:"al"}
				foo.Addresses = []Address{address1, address2}

				g := goson.Hash(foo, "Id")
				g.Array("Addresses", "City").Method("State")

				data, err := g.ToJson()
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`{"addresses":[{"city":"Ctr","state":"al"},{"city":"Tr","state":"al"}],"id":1}`))
			})
		})

		Context("Array", func() {
			It("Should include the nested array", func() {
				foo := User{Name: "Foo", Id: 1}
				address1 := Address{City:"Ctr", State:"al"}
				address2 := Address{City:"Tr", State:"al"}
				foo.Addresses = []Address{address1, address2}

				bar := User{Name: "Bar", Id: 2}
				address3 := Address{City:"Gr", State:"Gr"}
				address4 := Address{City:"Or", State:"al"}
				bar.Addresses = []Address{address3, address4}

				g := goson.Array([]User{foo, bar}, "Id")
				g.ArrayAlias("Addresses", "dir", "City", "State")

				data, err := g.ToJson()
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`[{"dir":[{"city":"Ctr","state":"al"},{"city":"Tr","state":"al"}],"id":1},{"dir":[{"city":"Gr","state":"Gr"},{"city":"Or","state":"al"}],"id":2}]`))
			})
		})
	})
})
