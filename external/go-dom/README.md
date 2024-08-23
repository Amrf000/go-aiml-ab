[![GoDoc](https://godoc.org/github.com/bserdar/go-dom?status.svg)](https://godoc.org/github.com/bserdar/go-dom)
[![Go Report Card](https://goreportcard.com/badge/github.com/bserdar/go-dom)](https://goreportcard.com/report/github.com/bserdar/go-dom)
[![Build Status](https://github.com/bserdar/go-dom/actions/workflows/CI.yml/badge.svg?branch=main)](https://github.com/bserdar/go-dom/actions/workflows/CI.yml)
# XML DOM Implementation for Go

This is a library to programmatically create, read, and modify XML
documents using Go without marshaling/unmarshaling them into
structures. It somewhat follows the XML DOM specification with the
following exceptions:

 * This implementation preserves and exposes XML namespace prefixes
 * Elements can be created with namespaces and prefixes
 * CDATA sections are converted to text nodes
 
## Namespace Normalization

The `Document` interface has a `NormalizeNamespaces()` method that
should be called before serializing the document. It does the following:

 * If there are elements/attributes with namespaces with no associated
   prefixes, then it create prefixes for them
 * If there are elements/attributes with prefixes that are not
   defined, it returns error
   
## Serialization

To parse XML documents, use the `Parse` function with an
`xml.Decoder.` By changing the `Strict` and `Entities` fields of a
`Decoder`, this parser can be used to parse HTML data.

To encode a `Document` as XML, first call `NormalizeNamespaces()`
function, and then use the `Encode` function.

