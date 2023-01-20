Multipart Upload: https://www.sobyte.net/post/2022-03/go-multipart-form-data/

```rpc

rpc_version = 1.0.0

name = example          # name of your backend app
version = v0.0.1        # version of your schema

enum Kind uint32 {
  _ = flag
  USER
  ADMIN
}

message Empty {

}

message User {
  id: uint64 {
    json = id,omitempty
    go.field.name = ID
    go.tag.db = id
  }

  username: string {
    json = USERNAME
    go.tag.db = username
  }

  role: string {
    json =
    go.tag.db = -
  }
}

message SearchFilter {
  q?: string
}

message ComplexType {
  meta: map<string,any>
  metaNestedExample: map<string,map<string,uint32>>
  namesList: []string
  numsList: []int64
  doubleArray: [][]string
  listOfMaps: []map<string,uint32>
  listOfUsers: []User
  mapOfUsers: map<string,User>
  user: User
}

message Base {
  id: string
}

message Child {
  ...Base
  firstname: string {
    json =                            # if nothing is placed, it considered as ignore
  }
}

service ExampleService {
  Ping()
  Status() => (status: bool)
  Version() => (version: Version)
  GetUser(header: map<string,string>, userID: uint64) => (code: uint32, user: User)
  FindUser(s: SearchFilter) => (name: string, user: User)

  findIDs(id: string) => (users: []User)

  UploadAvatar(files: MultiPart) => (status: Status)
  SubscribeEvent(type: string) => (events: stream Event)
}


```
