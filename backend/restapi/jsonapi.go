// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package restapi

// ResourceIdentifier represents a JSON:API resource identifier object:
// { "type": "documents", "id": "42" }
type ResourceIdentifier struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// ToOneRelationship is a JSON:API to-one relationship wrapper.
type ToOneRelationship struct {
	Data *ResourceIdentifier `json:"data"`
}

// ToManyRelationship is a JSON:API to-many relationship wrapper.
type ToManyRelationship struct {
	Data []ResourceIdentifier `json:"data"`
}

// PaginationMeta is optional, but handy if you want Ember to know counts.
type PaginationMeta struct {
	TotalItems int `json:"totalItems,omitempty"`
	PageSize   int `json:"pageSize,omitempty"`
	PageNumber int `json:"pageNumber,omitempty"`
}

// Links is optional JSON:API links object.
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}
