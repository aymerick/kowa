package commands

import (
	"fmt"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/server"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap Kowa",
	Long:  `Creates records in database`,
	Run:   bootstrap,
}

func bootstrap(cmd *cobra.Command, args []string) {
	// @todo Check that we are NOT in production

	// Insert oauth client
	oauthStorage := server.NewOAuthStorage()
	oauthStorage.SetupDefaultClient()

	db := models.NewDBSession()

	password, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		panic("Arg")
	}

	// Insert users
	userJeanClaude := models.User{
		Id:        "test",
		CreatedAt: time.Now(),
		Email:     "trucmush@wanadoo.fr",
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		Password:  string(password),
	}
	db.UsersCol().Insert(&userJeanClaude)

	userHenry := models.User{
		Id:        "hkanan",
		CreatedAt: time.Now(),
		Email:     "henrykanan@yahoo.com",
		FirstName: "Henry",
		LastName:  "Kanan",
		Password:  string(password),
	}
	db.UsersCol().Insert(&userHenry)

	// Insert sites
	siteJC1 := models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userJeanClaude.Id,
		CreatedAt:   time.Now(),
		Name:        "My site",
		Tagline:     "So powerfull !",
		Description: "You will be astonished by what my site is about",
	}
	db.SitesCol().Insert(&siteJC1)

	siteJC2 := models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userJeanClaude.Id,
		CreatedAt:   time.Now(),
		Name:        "My second site",
		Tagline:     "Very interesting",
		Description: "Our projects are so importants, please help us",
	}
	db.SitesCol().Insert(&siteJC2)

	siteH := models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userHenry.Id,
		CreatedAt:   time.Now(),
		Name:        "Ultimate petanque",
		Tagline:     "La petanque comme vous ne l'avez jamais vu",
		Description: "C'est vraiment le sport du futur. Messieurs, preparez vos boules !",
	}
	db.SitesCol().Insert(&siteH)

	// Insert posts
	var post models.Post

	for i := 1; i <= 30; i++ {
		post = models.Post{
			Id:          bson.NewObjectId(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			SiteId:      siteJC1.Id,
			PublishedAt: time.Now(),
			Title:       fmt.Sprintf("Post %d", i),
			Body: fmt.Sprintf(`Turba desolatas corpore iubet omnia patrios cingebant

## Dux quem fremida proque

Lorem markdownum noster tamen albida, **post retro**; haut hortator merito
imitabere? Carmina sua humanum adfata ponit ausus ore si inani natales inplevere
cupido sanior! Respiceret quam ad, vigil haec undam Ardea currus totis nec
diversa umbris esse!

1. Lac quas ille illi vidit nam armat
2. Et esse necem cultosque superis te molli
3. Ipsis inquit sonuere
4. Et idem
5. Raptas motu pro fuit solvitur
6. Verba interea miseri

## Turbae quo nova pugnae habeat

Violata est vidisse et sedula Melaneus *miseri*, nullius fertur. Venabula
candida magicaeque fere aestatem in quid, arma quam interdum, manifesta et alvum
[aequor petis](http://www.billmays.net/), Clymeneia. Ora Atridae illic. Deum
iram ergo in super lacrimis agros. In nutrix recessu Troius, tuis facientes
saevitiae aspicit potior.

## Quoque umerumque sepulcris equi

Incepta perpetuum balistave tandem reperitur tacuit libidine propago ecce quoque
serpit arduus. Navita est coniunx idem, penetrabile vivit, et de lux ira,
sollicitae haerebat. Sed cum miserantibus deus florente dentes ictu gutture
nosces. Exosa ense sic: ad, tibi.

## Inertes dierum nova

Oris populo, licet. Restabat videntur violenta loquebatur et glorior et flexus
Caras, modo istis.

- Omne fila oppugnant optari denique
- Novissima surgere
- Ait si soli gramina
- Ora fata acui dictos et missi

Cum lentus, et sinus, ferreus utroque est, non est promittes multamque summam
cognoram averso; fas. Totidem tectum flumina ubi prolem spiris in gurgite
**cadunt** copia suo hoc ebur, tua natae vellemque fumante stabulorum? Ictus
puerpera est o! Lycus supplice movere pars ecce obest cervina ulterius urna
medium tu forti tantumne altera excita moventur excita vestras! Thalamos qua
Iunonis non acumine nihil hamato ego tepidis et corvum.`, i),
		}
		db.PostsCol().Insert(&post)
	}

	post = models.Post{
		Id:          bson.NewObjectId(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		SiteId:      siteJC2.Id,
		PublishedAt: time.Now(),
		Title:       "This is a lonely",
		Body:        "It appears on my second website.",
	}
	db.PostsCol().Insert(&post)

	post = models.Post{
		Id:          bson.NewObjectId(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		SiteId:      siteH.Id,
		PublishedAt: time.Now(),
		Title:       "Hi, I am Henry",
		Body:        "Je me présente, je m'appelle Henry. Je voudrais bien réussir ma vie, être aimé. Être beau, gagner de l'argent. Puis surtout être intelligent. Mais pour tout ça il faudrait que je bosse à plein temps",
	}
	db.PostsCol().Insert(&post)

	// Insert events

	// Insert pages

	// Insert actions
}
