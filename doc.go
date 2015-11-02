// Copyright 2015 Kreditech Holding SSL GmbH & contributors. All rights reserved.
// Use of this source code is governed by a MIT-license.

/*
go-cbes - Golang ORM for Couchbase & ElastiSearch

go-cbes is a very fast ORM library for Golang that is using CouchBase and ElasticSearch as database.
It uses idiomatic Go to operate on databases, implementing struct to database mapping and acts as a lightweight
Go ORM framework. This library was designed to be supported by [Beego](http://beego.me/) or used as standalone
library as well to find a good balance between functionality and performance.

Requirements
 - ElasticSearch
 - Couchbase

Before using go-cbes make sure that you have installed and configure CouchBase and ElasticSearch. For CouchBase you
need to create your bucket manually, go-cbes will create automatically the ElasticSearch Index.

Look up the documentation on Github for more details (https://github.com/Kreditech/go-cbes).
*/
package cbes
