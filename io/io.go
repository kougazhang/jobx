package io

type Src interface{}

type Dst interface{}

type IO interface {
    Copy(src Src, dst Dst) error
}
