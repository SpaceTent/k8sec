// Copyright (c) 2026 SpaceTent
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package main

import (
    "fmt"
    
    v1 "k8s.io/api/core/v1"
)

type item struct {
    secret v1.Secret
}

func (i item) Title() string       { return i.secret.Name }
func (i item) Description() string { return fmt.Sprintf("Namespace: %s", i.secret.Namespace) }
func (i item) FilterValue() string { return i.secret.Name }
