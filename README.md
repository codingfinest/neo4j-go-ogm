# Neo4j-go-ogm - A Neo4j Object Graph Mapping Library for Golang runtime

The main goal of Neo4j-go-ogm is to get developers up and running quickly without writing unnecessary boiler plate code. Using Go Tags, Go runtime `structs` can be annotated and mapped to Neo4j domain entities. This project was hugely inspired by the Java version of Neo4j OGM.

## Quick start

```
go get -u github.com/codingfinest/neo4j-go-ogm
```

### Set up domain entities

The setup below declares 3 node entities (`Actor`, `Movie`, `Director`) and 1 relelationship entity (`Character`) which relates an `Actor` to a `Movie`. 

```
type Movie struct {
	gogm.Node `gogm:"label:FILM,label:PICTURE"`

	Title    string
	Released int64
	Tagline  string

	Characters []*Character `gogm:"direction:<-"`
	Directors  *Director  `gogm:"direction:<-"`
}

type Director struct {
	gogm.Node
	Name   string
	Movies []*Movie
}

type Actor struct {
	gogm.Node
	Name       string
	Characters []*Character
}

type Character struct {
	gogm.Relationship `gogm:"reltype:ACTED_IN"`

	Roles          []string
	Name           string
	NumberOfScenes int64
	Actor          *Actor `gogm:"startNode"`
	Movie          *Movie `gogm:"endNode"`
}
```

Relationships must be a pointer to a `struct` as can be seen in the `Movie`-`Director` relationship or a slice of pointer to `struct` as can be seen in `Characters` within `Actor` or `Movie`

### Persist/Load entities

```

	var config = &gogm.Config{
		"uri",
		"username",
		"password",
		gogm.NONE, /*log level*/
		false /*allow cyclic ref*/}

	var ogm = gogm.New(config)
	var session, err = ogm.NewSession(true)
	if err != nil {
		panic(err)
	}

	theMatrix := &Movie{}
	theMatrix.Title = "The Matrix"
	theMatrix.Released = 1999

	carrieAnne := &Actor{}
	carrieAnne.Name = "Carrie-Anne Moss"

	keanu := &Actor{}
	keanu.Name = "Keanu Reeves"

	carrieAnneMatrixCharacter := &Character{Movie: theMatrix, Actor: carrieAnne, Roles: []string{"Trinity"}}
	keanuReevesMatrixCharacter := &Character{Movie: theMatrix, Actor: keanu, Roles: []string{"Neo"}}

	theMatrix.Characters = append(theMatrix.Characters, carrieAnneMatrixCharacter, keanuReevesMatrixCharacter)

	//Persist the movie. This persists the actors as well and creates the relationships with their associated properties
	if err := session.Save(&theMatrix, nil); err != nil {
		panic(err)
	}

	var loadedMatrix *Movie
	if err := session.Load(&loadedMatrix, *theMatrix.ID, nil); err != nil {
		panic(err)
	}

	for _, character := range loadedMatrix.Characters {
		fmt.Println("Actor: " + character.Actor.Name)
	}
```

### Features
* **Save only deltas**: Persist only modified changes.
* **Node label inheritance**: Labels can be inherited from embedded node struct
* **Customizable node labels and relationship type**: Don't like the default node label or relationship type ? Easily customize them with the `label` and `reltype` struct tags respectively.
* **Runtime managed labels**: Dynamically manage your node labels at runtime
* **Transactions**: Commit or Rollback changes made to runtime objects
* **Rich Relationships**: Add properties to relationships
* **Persistence event**: Intercept events during the lifecyle of a runtime object. 
* **Custom queries**: Create custom queries to polulate runtime objects

### Struct Tags
* `id`: Entity identifier. Only primitive types are supported. Unique constraint is created on this field. 
* `unique`: Creates a unique constraint on this field.
* `index`: Creates an index on this field.
* `label`: Used to customize the node labels when tagged on the embedded `gogm.Node` `struct` or any embedded annonymous `struct` embedding `gogm.Node`. When used on a field with type `[]string`, it identifies that field as the source of runtime manage labels.
* `reltype`: Used to customize relationship type. It should be tagged on `gogm.Relationship` for relationship entities or fields within nodes that relate to other nodes.
* `name`: Denotes the name to use in the database for the property associated with this field.
* `direction`: Indicates the direction of relationship. Possible values are `<-` for incoming, `--` for undirected and `->` for outgoing. When a direction isn't specified, default is `->`.
* `startNode`: Denotes the start node of a relationship
* `endNode`: Denotes the end node of a relationship
* `-`: Ignore field




### Reporting bugs

Thanks for trying out the package! Any bug found should be documented with specific reproducible conditions on the githib issue page.



### Coming soon
* **Load options** Filters, Sort, Pagination etc

### LICENSE

MIT













