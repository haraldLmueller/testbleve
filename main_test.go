/*
Copyright © 2024 Harald Müller

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/stretchr/testify/assert"
)

var catalogPath = "index.bleve"

func Test_Bleve(t *testing.T) {
	os.RemoveAll(catalogPath)
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewUsing(catalogPath, indexMapping, bleve.Config.DefaultIndexType, bleve.Config.DefaultKVStore, nil)
	assert.NoError(t, err)
	testperson := insertToIndex(index, nil)
	q := bleve.NewQueryStringQuery("+" + testperson.Name.First + " +" + testperson.Name.Last)
	req := bleve.NewSearchRequestOptions(q, 100, 0, false)
	res, err := index.Search(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Hits))
	assert.Equal(t, "person50", res.Hits[0].ID)
	index.Close()
}

func Test_BleveBatch(t *testing.T) {
	os.RemoveAll(catalogPath)
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewUsing(catalogPath, indexMapping, bleve.Config.DefaultIndexType, bleve.Config.DefaultKVStore, nil)
	assert.NoError(t, err)
	batch := index.NewBatch()
	testperson := insertToIndex(nil, batch)
	err = index.Batch(batch)
	assert.NoError(t, err)
	q := bleve.NewQueryStringQuery("+" + testperson.Name.First + " +" + testperson.Name.Last)
	req := bleve.NewSearchRequestOptions(q, 100, 0, false)
	res, err := index.Search(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Hits))
	assert.Equal(t, "person50", res.Hits[0].ID)
	index.Close()
}

func Test_BleveWithDelay(t *testing.T) {
	os.RemoveAll(catalogPath)
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewUsing(catalogPath, indexMapping, bleve.Config.DefaultIndexType, bleve.Config.DefaultKVStore, nil)
	assert.NoError(t, err)
	testperson := insertToIndex(index, nil)
	q := bleve.NewQueryStringQuery("+" + testperson.Name.First + " +" + testperson.Name.Last)
	req := bleve.NewSearchRequestOptions(q, 100, 0, false)
	time.Sleep(time.Second * 40)
	res, err := index.Search(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Hits))
	assert.Equal(t, "person50", res.Hits[0].ID)
	index.Close()
}

func insertToIndex(index bleve.Index, batch *bleve.Batch) Person {
	var testperson Person
	myRand := rand.New(rand.NewSource(123))
	for i := range 5000 {
		p := createRandomPerson(int64(i), 500, myRand)
		if batch != nil {
			batch.Index("person"+strconv.Itoa(i), p)
		} else {
			index.Index("person"+strconv.Itoa(i), p)
		}
		if i == 50 {
			testperson = *p
			jsonByte, _ := json.Marshal(p)
			filename := "jsonfiles/person_" + strconv.Itoa(i) + ".json"
			err := os.MkdirAll("jsonfiles", 0777)
			checkFatal(err, "create jsonfiles directory")
			err = os.WriteFile(filename, jsonByte, 0666)
			checkFatal(err, "create test person")
		}
	}
	return testperson
}

type Name struct {
	First string
	Last  string
}

type Person struct {
	ID      int64
	Name    Name
	Age     int
	City    string
	Like    []int64
	Comment map[string]string
}

func createRandomPerson(id int64, maxID int64, myRand *rand.Rand) *Person {
	p := Person{ID: int64(id),
		Name: Name{First: names[myRand.Intn(len(names))], Last: lastNames[myRand.Intn(len(lastNames))]},
		City: cities[myRand.Intn(len(cities))],
		Age:  myRand.Intn(103),
	}
	var friends []int64
	for range myRand.Intn(3) + 1 {
		friends = append(friends, myRand.Int63n(maxID))
	}
	p.Like = friends
	p.Comment = make(map[string]string, 0)
	for range myRand.Intn(6) + 350 {
		p.Comment[names[myRand.Intn(len(names))]] = cities[myRand.Intn(len(cities))]
	}

	return &p
}

func checkFatal(err error, txt string) {
	if err != nil {
		slog.Default().Error(txt, "error", err)
		os.Exit(1)
	}
}

var names = []string{"James", "Mary", "Robert", "Patricia", "John", "Jennifer", "Michael", "Linda", "David", "Elizabeth", "William", "Barbara", "Richard",
	"Susan", "Joseph", "Jessica", "Thomas", "Sarah", "Charles", "Karen", "Christopher", "Lisa", "Daniel", "Nancy", "Matthew", "Betty", "Anthony",
	"Margaret", "Mark", "Sandra", "Donald", "Ashley", "Steven", "Kimberly", "Paul", "Emily", "Andrew", "Donna", "Joshua", "Michelle", "Kenneth",
	"Carol", "Kevin", "Amanda", "Brian", "Dorothy", "George", "Melissa", "Timothy", "Deborah", "Ronald", "Stephanie", "Edward", "Rebecca", "Jason",
	"Sharon", "Jeffrey", "Laura", "Ryan", "Cynthia", "Jacob", "Kathleen", "Gary", "Amy", "Nicholas", "Angela", "Eric", "Shirley", "Jonathan",
	"Anna", "Stephen", "Brenda", "Larry", "Pamela", "Justin", "Emma", "Scott", "Nicole", "Brandon", "Helen", "Benjamin", "Samantha", "Samuel",
	"Katherine", "Gregory", "Christine", "Alexander", "Debra", "Frank", "Rachel", "Patrick", "Carolyn", "Raymond", "Janet", "Jack", "Catherine",
	"Dennis", "Maria", "Jerry", "Heather", "Tyler", "Diane", "Aaron", "Ruth", "Jose", "Julie", "Adam", "Olivia", "Nathan", "Joyce", "Henry",
	"Virginia", "Douglas", "Victoria", "Zachary", "Kelly", "Peter", "Lauren", "Kyle", "Christina", "Ethan", "Joan", "Walter", "Evelyn", "Noah",
	"Judith", "Jeremy", "Megan", "Christian", "Andrea", "Keith", "Cheryl", "Roger", "Hannah", "Terry", "Jacqueline", "Gerald", "Martha", "Harold",
	"Gloria", "Sean", "Teresa", "Austin", "Ann", "Carl", "Sara", "Arthur", "Madison", "Lawrence", "Frances", "Dylan", "Kathryn", "Jesse",
	"Janice", "Jordan", "Jean", "Bryan", "Abigail", "Billy", "Alice", "Joe", "Julia", "Bruce", "Judy", "Gabriel", "Sophia", "Logan", "Grace",
	"Albert", "Denise", "Willie", "Amber", "Alan", "Doris", "Juan", "Marilyn", "Wayne", "Danielle", "Elijah", "Beverly", "Randy", "Isabella",
	"Roy", "Theresa", "Vincent", "Diana", "Ralph", "Natalie", "Eugene", "Brittany", "Russell", "Charlotte", "Bobby", "Marie", "Mason", "Kayla",
	"Philip", "Alexis", "Louis", "Lori"}
var lastNames = []string{"Abraham", "Allan", "Alsop", "Anderson", "Arnold", "Avery", "Bailey", "Baker", "Ball", "Bell", "Berry", "Black", "Blake", "Bond",
	"Bower", "Brown", "Buckland", "Burgess", "Butler", "Cameron", "Campbell", "Carr", "Chapman", "Churchill", "Clark", "Clarkson", "Coleman", "Cornish",
	"Davidson", "Davies", "Dickens", "Dowd", "Duncan", "Dyer", "Edmunds", "Ellison", "Ferguson", "Fisher", "Forsyth", "Fraser", "Gibson", "Gill", "Glover",
	"Graham", "Grant", "Gray", "Greene", "Hamilton", "Hardacre", "Harris", "Hart", "Hemmings", "Henderson", "Hill", "Hodges", "Howard", "Hudson", "Hughes",
	"Hunter", "Ince", "Jackson", "James", "Johnston", "Jones", "Kelly", "Kerr", "King", "Knox", "Lambert", "Langdon", "Lawrence", "Lee", "Lewis", "Lyman",
	"MacDonald", "Mackay", "Mackenzie", "MacLeod", "Manning", "Marshall", "Martin", "Mathis", "May", "McDonald", "McLean", "McGrath", "Metcalfe", "Miller",
	"Mills", "Mitchell", "Morgan", "Morrison", "Murray", "Nash", "Newman", "Nolan", "North", "Ogden", "Oliver", "Paige", "Parr", "Parsons", "Paterson",
	"Payne", "Peake", "Peters", "Piper", "Poole", "Powell", "Pullman", "Quinn", "Rampling", "Randall", "Rees", "Reid", "Roberts", "Robertson", "Ross", "Russell",
	"Rutherford", "Sanderson", "Scott", "Sharp", "Short", "Simpson", "Skinner", "Slater", "Smith", "Springer", "Stewart", "Sutherland", "Taylor", "Terry",
	"Thomson", "Tucker", "Turner", "Underwood", "Vance", "Vaughan", "Walker", "Wallace", "Walsh", "Watson", "Welch", "White", "Wilkins", "Wilson",
	"Wright", "Young"}
var cities = []string{"Bladensburg", "Brambleton", "Edenburg", "Dubois", "Cotopaxi", "Sperryville", "Alleghenyville", "Westboro", "Tonopah", "Fowlerville",
	"Venice", "Wanship", "Diaperville", "Haring", "Morriston", "Kenvil", "Dahlen", "Canby", "Basye", "Marienthal", "Sutton", "Elwood",
	"Tilleda", "Crenshaw", "Loveland", "Canoochee", "Newkirk", "National", "Chesterfield", "Draper", "Turah", "Hall", "Dragoon", "Summertown", "Sims",
	"Guthrie", "Vivian", "Tuttle", "Ladera", "Drummond", "Ezel", "Marne", "Lookingglass", "Shasta", "Vandiver", "Sharon", "Glendale", "Loomis",
	"Statenville", "Gouglersville", "Sehili", "Catherine", "Whitmer", "Grimsley", "Salix", "Kersey", "Springdale", "Thermal", "Witmer", "Virgie",
	"Wakulla", "Indio", "Unionville", "Loretto", "Sabillasville", "Gracey", "Blodgett", "Aguila", "Harleigh", "Avalon", "Fairview",
	"Esmont", "Cascades", "Cleary", "Reno", "Holtville", "Lumberton", "Keller", "Caspar", "Biddle", "Dexter", "Whitehaven", "Fidelis", "Drytown",
	"Dorneyville", "Rivereno", "Independence", "Bodega", "Wanamie", "Townsend", "Caron", "Guilford", "Gallina", "Manila", "Itmann", "Whitewater",
	"Templeton", "Jessie", "Sena", "Charco", "Jamestown", "Imperial", "Vincent", "Nelson", "Abrams", "Glasgow", "Lynn", "Sugartown", "Navarre",
	"Marion", "Sanders", "Spelter", "Santel", "Outlook", "Ypsilanti", "Dotsero", "Mathews", "Loyalhanna", "Libertytown", "Terlingua", "Hackneyville",
	"Driftwood", "Stockdale", "Bynum", "Harrison", "Morningside", "Churchill", "Gambrills", "Brule", "Fairhaven", "Hinsdale", "Babb", "Buxton",
	"Biehle", "Catharine", "Dunbar", "Klagetoh", "Blandburg", "Roberts", "Romeville", "Hachita", "Leming", "Saranap", "Elliott", "Ronco", "Rossmore",
	"Bowie", "Roderfield", "Devon", "Trucksville", "Ribera", "Watchtower", "Orason", "Haena", "Fruitdale", "Riceville", "Urbana", "Moscow",
	"Fulford", "Cassel", "Shawmut", "Corinne", "Edmund", "Naomi", "Clara", "Duryea", "Chloride", "Axis", "Villarreal", "Talpa", "Rodman", "Goochland",
	"Deercroft", "Jacksonburg", "Kanauga", "Springville", "Concho", "Matheny", "Temperanceville", "Salunga", "Elfrida", "Stollings", "Lindisfarne",
	"Kimmell", "Fillmore", "Belmont", "Mansfield", "Fairforest", "Finzel", "Shelby", "Brenton", "Fairlee", "Brownlee", "Yettem", "Richmond", "Jeff",
	"Umapine", "Cuylerville", "Carbonville", "Alamo"}
