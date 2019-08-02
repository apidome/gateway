# caf
Cloud Application Firewall

## Project Structure Source
https://github.com/golang-standards/project-layout


## High Level Design
### Proxy
- Listen to requests from client
- Forwards request as-is to the target

### REST Engine
- Validate REST schema (Implement RFC)
- Check POST/PUT/PATCH/etc (Methods that contain a JSON body)
- Validate JSON structure (Make sure all required fields exist)
- Validate the data of each field (RegEx at this point, Validation functions in the future)
