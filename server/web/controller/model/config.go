/**
 *
 * (c) Copyright Ascensio System SIA 2025
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package model

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type Config struct {
	jwt.RegisteredClaims
	Document     Document     `json:"document"`
	DocumentType string       `json:"documentType"`
	EditorConfig EditorConfig `json:"editorConfig"`
	Type         string       `json:"type"`
	Token        string       `json:"token,omitempty"`
}

type Document struct {
	FileType    string      `json:"fileType"`
	Key         string      `json:"key"`
	Title       string      `json:"title"`
	URL         string      `json:"url"`
	Permissions Permissions `json:"permissions"`
}

type EditorConfig struct {
	User          User          `json:"user"`
	CallbackURL   string        `json:"callbackUrl"`
	Customization Customization `json:"customization,omitempty"`
	Lang          string        `json:"lang,omitempty"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Customization struct {
	Goback Goback `json:"goback"`
}

type Goback struct {
	RequestClose bool `json:"requestClose"`
}
