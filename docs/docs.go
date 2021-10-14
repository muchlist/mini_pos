// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "whois.muchlis@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "description": "login menggunakan userID dan password untuk mendapatkan JWT Token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Access"
                ],
                "summary": "login",
                "operationId": "user-login",
                "parameters": [
                    {
                        "description": "Body raw JSON",
                        "name": "ReqBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UserLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.UserLoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample400"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample500"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/profile": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "menampilkan profile berdasarkan user yang login saat ini",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Access"
                ],
                "summary": "get current profile",
                "operationId": "user-profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.UserModel"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/refresh": {
            "post": {
                "description": "mendapatkan token dengan tambahan waktu expired menggunakan refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Access"
                ],
                "summary": "refresh token",
                "operationId": "user-refresh",
                "parameters": [
                    {
                        "description": "Body raw JSON",
                        "name": "ReqBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UserRefreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.UserRefreshTokenResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample400"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample500"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "menampilkan daftar user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Access"
                ],
                "summary": "find user",
                "operationId": "user-find",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset cursor untuk skip data sebanyak offsite",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search apabila di isi akan melakukan pencarian berdasarkan nama",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/dto.UserModel"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample400"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample500"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "menambahkan user pada merchant sesuai usr owner, endpoint ini membutuhkan hak akses owner, sedangkan akun owner dapat didapatkan ketika membuat Merchant",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Access"
                ],
                "summary": "register user",
                "operationId": "user-register",
                "parameters": [
                    {
                        "description": "Body raw JSON",
                        "name": "ReqBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UserRegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/wrap.RespMsgExample"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample400"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample500"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "menampilkan user berdasarkan userID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Access"
                ],
                "summary": "get user by ID",
                "operationId": "user-get",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
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
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.UserModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample400"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample500"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "melakukan perubahan data pada user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Access"
                ],
                "summary": "edit user",
                "operationId": "user-edit",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Body raw JSON",
                        "name": "ReqBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UserEditRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.UserModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample400"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample500"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "menghapus user berdasarkan userID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Access"
                ],
                "summary": "delete user by ID",
                "operationId": "user-delete",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/wrap.RespMsgExample"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample400"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/wrap.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "error": {
                                            "$ref": "#/definitions/wrap.ErrorExample500"
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
        "dto.UserEditRequest": {
            "type": "object",
            "properties": {
                "def_outlet": {
                    "type": "integer",
                    "example": 1
                },
                "email": {
                    "type": "string",
                    "example": "example@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "muchlis"
                },
                "role": {
                    "type": "string",
                    "example": "employee"
                }
            }
        },
        "dto.UserLoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "dto.UserLoginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"
                },
                "def_outlet": {
                    "type": "integer",
                    "example": 1
                },
                "email": {
                    "type": "string",
                    "example": "example@example.com"
                },
                "expired": {
                    "type": "integer",
                    "example": 1631341964
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "merchant_id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "muchlis"
                },
                "refresh_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"
                },
                "role": {
                    "type": "string",
                    "example": "owner,employee"
                }
            }
        },
        "dto.UserModel": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer",
                    "example": 1631341964
                },
                "def_outlet": {
                    "type": "integer",
                    "example": 1
                },
                "email": {
                    "type": "string",
                    "example": "example@example.com"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "merchant_id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "muchlis"
                },
                "role": {
                    "type": "string",
                    "example": "owner,employee"
                },
                "updated_at": {
                    "type": "integer",
                    "example": 1631341964
                }
            }
        },
        "dto.UserRefreshTokenRequest": {
            "type": "object",
            "properties": {
                "refresh_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"
                }
            }
        },
        "dto.UserRefreshTokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"
                },
                "expired": {
                    "type": "integer",
                    "example": 1631341964
                }
            }
        },
        "dto.UserRegisterRequest": {
            "type": "object",
            "properties": {
                "def_outlet": {
                    "type": "integer",
                    "example": 1
                },
                "email": {
                    "type": "string",
                    "example": "example@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "muchlis"
                },
                "password": {
                    "type": "string",
                    "example": "password123"
                },
                "role": {
                    "type": "string",
                    "example": "owner,employee"
                }
            }
        },
        "wrap.ErrorExample400": {
            "type": "object",
            "properties": {
                "causes": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "causes 1",
                        "causes 2"
                    ]
                },
                "error": {
                    "type": "string",
                    "example": "unauthorized"
                },
                "message": {
                    "type": "string",
                    "example": "Unauthorized, memerlukan hak akses [ADMIN]"
                },
                "status": {
                    "type": "integer",
                    "example": 401
                }
            }
        },
        "wrap.ErrorExample500": {
            "type": "object",
            "properties": {
                "causes": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "ERROR: argument of WHERE must be type boolean. not type integer (SQLSTATE 42804)"
                    ]
                },
                "error": {
                    "type": "string",
                    "example": "internal_server_error"
                },
                "message": {
                    "type": "string",
                    "example": "gagal saat penghapusan item"
                },
                "status": {
                    "type": "integer",
                    "example": 500
                }
            }
        },
        "wrap.Resp": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "x-nullable": true
                },
                "error": {
                    "type": "object",
                    "x-nullable": true
                }
            }
        },
        "wrap.RespMsgExample": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string",
                    "x-nullable": true,
                    "example": "Data dengan ID xxx berhasil di [Create/Delete]"
                },
                "error": {
                    "type": "object",
                    "x-nullable": true
                }
            }
        }
    },
    "securityDefinitions": {
        "bearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:3500",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "mini_pos API",
	Description: "Mini Pos Api",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
