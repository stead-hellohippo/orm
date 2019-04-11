// Code generated by prana; DO NOT EDIT.

// Package database_tests contains a tests for database repository
// Auto-generated at Thu, 11 Apr 2019 16:37:34 CEST

package database_test

import (
	"database/sql"

	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/example/database"
	"github.com/phogolabs/orm/example/database/model"
	"github.com/phogolabs/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserRepository", func() {
	var (
		repository *database.UserRepository
		entity     *model.User
	)

	BeforeEach(func() {
		repository = &database.UserRepository{
			Gateway: gateway,
		}

		entity = &model.User{
			ID:        1,
			FirstName: "John",
			LastName:  schema.NullStringFrom("Doe"),
		}
	})

	AfterEach(func() {
		_, err := gateway.Exec(orm.SQL("DELETE FROM users"))
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("SelectAll", func() {
		It("returns no records", func() {
			records, err := repository.SelectAll()
			Expect(err).NotTo(HaveOccurred())
			Expect(records).To(BeEmpty())
		})

		Context("when there is a record", func() {
			BeforeEach(func() {
				Expect(repository.Insert(entity)).To(Succeed())
			})

			It("returns all records", func() {
				records, err := repository.SelectAll()
				Expect(err).NotTo(HaveOccurred())
				Expect(records).To(HaveLen(1))
				Expect(records[0]).To(Equal(entity))
			})
		})
	})

	Describe("SelectByPK", func() {
		It("return a record by primary key", func() {
			Expect(repository.Insert(entity)).To(Succeed())

			record, err := repository.SelectByPK(entity.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(record).To(Equal(entity))
		})

		Context("when the record does not exist", func() {
			It("returns an error", func() {
				record, err := repository.SelectByPK(entity.ID)
				Expect(err).To(Equal(sql.ErrNoRows))
				Expect(record).To(BeNil())
			})
		})
	})

	Describe("SearchAll", func() {
		It("returns all records for given RQL", func() {
			Expect(repository.Insert(entity)).To(Succeed())

			records, err := repository.SearchAll(&orm.RQLQuery{})
			Expect(err).NotTo(HaveOccurred())
			Expect(records).To(HaveLen(1))
			Expect(records[0]).To(Equal(entity))
		})
	})

	Describe("Insert", func() {
		It("inserts a new member successfully", func() {
			Expect(repository.Insert(entity)).To(Succeed())

			record, err := repository.SelectByPK(entity.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(record).To(Equal(entity))
		})
	})

	Describe("UpdateByPK", func() {
		BeforeEach(func() {
			Expect(repository.Insert(entity)).To(Succeed())
		})

		It("updates a record by primary key", func() {
			Expect(repository.UpdateByPK(entity)).To(Succeed())

			record, err := repository.SelectByPK(entity.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(record).To(Equal(entity))
		})
	})

	Describe("DeleteByPK", func() {
		BeforeEach(func() {
			Expect(repository.Insert(entity)).To(Succeed())
		})

		It("deletes a record by primary key", func() {
			Expect(repository.DeleteByPK(entity.ID)).To(Succeed())

			record, err := repository.SelectByPK(entity.ID)
			Expect(err).To(Equal(sql.ErrNoRows))
			Expect(record).To(BeNil())
		})
	})
})
