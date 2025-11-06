package httpserver

// @title FamilyTree API
// @version 1.0
// @description To get a token manually:
// @description ```
// @description curl -X POST 'http://localhost:15101/realms/familytree/protocol/openid-connect/token' \
// @description   -H 'Content-Type: application/x-www-form-urlencoded' \
// @description   -d 'client_id=familytree-cli' \
// @description   -d 'username=testuser' \
// @description   -d 'password=S3cr3t!' \
// @description   -d 'grant_type=password'
// @description ```
// @description
// @description Extract and use the access_token:
// @description ```
// @description TOKEN=$(curl -s -X POST 'http://localhost:15101/realms/familytree/protocol/openid-connect/token' \
// @description   -H 'Content-Type: application/x-www-form-urlencoded' \
// @description   -d 'client_id=familytree-cli' \
// @description   -d 'username=testuser' \
// @description   -d 'password=S3cr3t!' \
// @description   -d 'grant_type=password' | jq -r '.access_token')
// @description
// @description curl -H "Authorization: Bearer $TOKEN" http://localhost:15000/v1/persons
// @description ```

// @contact.name FamilyTree Support
// @contact.url https://github.com/rogerwesterbo/familytree
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url https://github.com/rogerwesterbo/familytree/blob/main/LICENSE

// @host localhost:15000
// @BasePath /

// @schemes http https
// @produce json
// @consumes json

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl http://localhost:15101/realms/familytree/protocol/openid-connect/token

// @tag.name Persons
// @tag.description Operations related to persons in the family tree

// @tag.name Relationships
// @tag.description Operations related to relationships between persons
