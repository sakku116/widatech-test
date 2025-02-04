// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/invoices": {
            "get": {
                "description": "set date_from and date_to to get profit_total and cash_transaction_total.\\nif date_from and date_to is set, pagination will be ignored for profit calculation.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "get invoice list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "DD-MM-YYYY, fill up date_from and date_to to get profit_total and cash_transaction_total. WARNING: if set, pagination will be ignored for profit calculation.",
                        "name": "date_from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "DD-MM-YYYY, fill up date_from and date_to to get profit_total and cash_transaction_total. WARNING: if set, pagination will be ignored for profit calculation.",
                        "name": "date_to",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "CASH",
                            "CREDIT"
                        ],
                        "type": "string",
                        "x-enum-varnames": [
                            "InvoicePaymentType_CASH",
                            "InvoicePaymentType_CREDIT"
                        ],
                        "description": "leave empty to query all payment types",
                        "name": "payment_type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "query",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "invoice_no",
                            "customer_name",
                            "sales_person_name"
                        ],
                        "type": "string",
                        "description": "leave empty to query by all queriable fields",
                        "name": "query_by",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "created_at",
                            "updated_at",
                            "invoice_no"
                        ],
                        "type": "string",
                        "default": "updated_at",
                        "name": "sort_by",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "asc",
                            "desc"
                        ],
                        "type": "string",
                        "default": "desc",
                        "x-enum-varnames": [
                            "SortOrder_asc",
                            "SortOrder_desc"
                        ],
                        "name": "sort_order",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.BaseJSONResp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.GetInvoiceListRespData"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "create new invoice",
                "parameters": [
                    {
                        "description": "create invoice request",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateInvoiceReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.BaseJSONResp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.CreateInvoiceResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/invoices/import": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "import invoice from xlsx",
                "parameters": [
                    {
                        "type": "file",
                        "description": "xlsx file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.BaseJSONResp"
                        }
                    }
                }
            }
        },
        "/invoices/no/{invoice_no}": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "delete invoice by invoice no",
                "parameters": [
                    {
                        "type": "string",
                        "description": "invoice no",
                        "name": "invoice_no",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.BaseJSONResp"
                        }
                    }
                }
            }
        },
        "/invoices/{invoice_uuid}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "get invoice detail",
                "parameters": [
                    {
                        "type": "string",
                        "description": "invoice uuid",
                        "name": "invoice_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.BaseJSONResp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.GetInvoiceDetailRespData"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "delete invoice",
                "parameters": [
                    {
                        "type": "string",
                        "description": "invoice uuid",
                        "name": "invoice_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.BaseJSONResp"
                        }
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "update invoice",
                "parameters": [
                    {
                        "type": "string",
                        "description": "invoice uuid",
                        "name": "invoice_uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update invoice request",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UpdateInvoiceReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.BaseJSONResp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.UpdateInvoiceRespData"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.BaseJSONResp": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "detail": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.CreateInvoiceReq": {
            "type": "object",
            "required": [
                "customer_name",
                "payment_type",
                "product_uuids",
                "sales_person_name"
            ],
            "properties": {
                "customer_name": {
                    "type": "string"
                },
                "notes": {
                    "type": "string"
                },
                "payment_type": {
                    "$ref": "#/definitions/enum.InvoicePaymentType"
                },
                "product_uuids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "sales_person_name": {
                    "type": "string"
                }
            }
        },
        "dto.CreateInvoiceResp": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "customer_name": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "invoice_no": {
                    "type": "string"
                },
                "notes": {
                    "type": "string"
                },
                "payment_type": {
                    "$ref": "#/definitions/enum.InvoicePaymentType"
                },
                "sales_person_name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "dto.GetInvoiceDetailRespData": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "customer_name": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "invoice_no": {
                    "type": "string"
                },
                "notes": {
                    "type": "string"
                },
                "payment_type": {
                    "$ref": "#/definitions/enum.InvoicePaymentType"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.BaseProductResp"
                    }
                },
                "sales_person_name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "dto.GetInvoiceListRespData": {
            "type": "object",
            "properties": {
                "cash_transaction_total": {
                    "type": "integer"
                },
                "current_page": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.GetInvoiceListRespData_DataItem"
                    }
                },
                "profit_total": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                },
                "total_page": {
                    "type": "integer"
                }
            }
        },
        "dto.GetInvoiceListRespData_DataItem": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "customer_name": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "invoice_no": {
                    "type": "string"
                },
                "notes": {
                    "type": "string"
                },
                "payment_type": {
                    "$ref": "#/definitions/enum.InvoicePaymentType"
                },
                "product_total": {
                    "type": "integer"
                },
                "sales_person_name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "dto.UpdateInvoiceReq": {
            "type": "object",
            "properties": {
                "customer_name": {
                    "type": "string"
                },
                "notes": {
                    "description": "use 'null' to set explicitly to null",
                    "type": "string"
                },
                "payment_type": {
                    "$ref": "#/definitions/enum.InvoicePaymentType"
                },
                "product_uuids": {
                    "description": "use empty array [] remove all products",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "sales_person_name": {
                    "type": "string"
                }
            }
        },
        "dto.UpdateInvoiceRespData": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "customer_name": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "invoice_no": {
                    "type": "string"
                },
                "notes": {
                    "type": "string"
                },
                "payment_type": {
                    "$ref": "#/definitions/enum.InvoicePaymentType"
                },
                "sales_person_name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "enum.InvoicePaymentType": {
            "type": "string",
            "enum": [
                "CASH",
                "CREDIT"
            ],
            "x-enum-varnames": [
                "InvoicePaymentType_CASH",
                "InvoicePaymentType_CREDIT"
            ]
        },
        "enum.SortOrder": {
            "type": "string",
            "enum": [
                "asc",
                "desc"
            ],
            "x-enum-varnames": [
                "SortOrder_asc",
                "SortOrder_desc"
            ]
        },
        "model.BaseProductResp": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "invoice_uuid": {
                    "type": "string"
                },
                "item_name": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                },
                "total_cost_of_goods_sold": {
                    "type": "integer"
                },
                "total_price_sold": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Widatech Test",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
