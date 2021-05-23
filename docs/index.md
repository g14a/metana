Metana is a cli abstract migration tool written in Go for Go services. It is inspired by real frustrations of not being able to track/store and run migrations in Go.

The motivation however, behind creating this tool, is to abstract away the database part. If your task can be completed with Pure Go or via a Go driver of your service, then this is for you. Since it makes use of the Go runtime, you can even perform database migrations like PostgreSQL, Mongo, Redis, Elasticsearch, GCP Buckets etc. You just need to be able to interact with your data store or complete your task using Go.

The main use case is when you won't be able to do everything with SQL or No-SQL syntax. There might be some tasks where you need to aggregate data, iterate over them, and do business related stuff with the retrieved data. All you need to know is Go syntax and write a Go program.

![OpenSource](https://img.shields.io/badge/Open%20Source-000000?style=for-the-badge&logo=github)
![go](https://img.shields.io/badge/-Written%20In%20Go-00add8?style=for-the-badge&logo=Go&logoColor=ffffff)
![cli](https://img.shields.io/badge/-Build%20for%20CLI-000000?style=for-the-badge&logo=Powershell&logoColor=ffffff)

[![Go Report Card](https://goreportcard.com/badge/github.com/g14a/metana)](https://goreportcard.com/report/github.com/g14a/metana)
[![Go Workflow Status](https://github.com/g14a/metana/workflows/Go/badge.svg)](https://github.com/g14a/metana/workflows/Go/badge.svg)

Checkout my blog at https://g14a.github.io