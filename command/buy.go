package command

import (
	"fmt"
	"github.com/magazin/command/dto"
	"github.com/magazin/data"
	"github.com/magazin/data/models"
	"github.com/magazin/data/repository"
	"github.com/magazin/services"
	"regexp"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

var buyHelpTemplate = `Name:
	{{.HelpName}} - {{.Usage}}

Usage:
	{{.HelpName}} [options] source

Options:
	{{range .VisibleFlags}}{{.}}
	{{end}}
Examples:
	1. Buy 10 of product with id of prodID1
		 > magazin {{.HelpName}} --uid=987 prodID1@10
	2. Buy 15 of prodID1 and one prodID2
		> magazin {{.HelpName}} --uid=988 prodID2 prodID1@15

`

const (
	phoneNotification = iota
	emailNotification
)

var buyCommand = &cli.Command{
	Name:               "buy",
	HelpName:           "buy",
	Usage:              "buy products",
	CustomHelpTemplate: buyHelpTemplate,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "user-id",
			Aliases: []string{"uid"},
			Usage:    "user id",
			Required: true,
		},
		&cli.StringFlag{
			Name: "notification-mode",
			Aliases: []string{"notif"},
			Usage: "set mode of notification: (phone, email)",
			Required: true,
		},
	},
	Before: func(c *cli.Context) error {
		if _, err := data.Begin(false).User(c.String("uid")); err != nil {
			return err
		}

		return validateBuyCommand(c)
	},
	Action: func(c *cli.Context) (err error) {

		return Buy{
			uid:       c.String("uid"),
			purchases: productsPurchased(c),
		}.Run()
	},
}

// Buy holds buy operation state.
type Buy struct {
	uid       string
	purchases []dto.Purchase
	notification
}

func (b Buy) Run() error {
	user, err := data.Begin(false).User(b.uid)
	if err != nil {
		return err
	}
	purchasedProds := make([]dto.ProductPurchased, 0)
	transact := data.Begin(true)

	for _, purchase := range b.purchases {
		purchasedProd, err := doPurchase(transact, purchase)
		if err != nil {
			return err
		}
		purchasedProds = append(purchasedProds, purchasedProd)
	}
	transact.Commit()

	result, err := sendNotification(user, purchasedProds, b.mode)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func sendNotification(
	user *models.User,
	products []dto.ProductPurchased,
	mode int) (string, error) {

	var sb strings.Builder
	sb.WriteString("Dear User,\nYou have purchased the following items:\n\n")
	var tcost float64
	for _, product := range products {
		tcost += product.Price * float64(product.Quantity)
		sb.WriteString(fmt.Sprintf("\t%s, %d many of them, each %v somoni\n", product.Name, product.Quantity, product.Price))
	}
	sb.WriteString(fmt.Sprintf("\ntotal cost: %.2f somoni", tcost))

	message := sb.String()

	if mode == emailNotification {
		email := services.Email{}
		if err := email.Send(user.Email, message); err != nil {
			return "", err
		}
		return message, nil
	}
	sms := services.SMS{}
	if err := sms.Send(user.Phone, message); err != nil {
		return "", err
	}
	return message, nil
}

func doPurchase(repo *repository.Transact, purchase dto.Purchase) (dto.ProductPurchased, error) {
	if err := decreaseQuant(repo, purchase); err != nil {
		return dto.ProductPurchased{}, err
	}
	product, err := repo.Product(purchase.ProductID)
	if err != nil {
		return dto.ProductPurchased{}, err
	}
	return dto.ProductPurchased{
		Product: dto.Product{
			Name:  product.Name,
			Price: product.Price,
		},
		Purchase: dto.Purchase{
			ProductID: product.ID,
			Quantity:  purchase.Quantity,
		},
	}, nil
}

func decreaseQuant(repo *repository.Transact, purchase dto.Purchase) error {
	quant, err := repo.Quantity(purchase.ProductID)
	if err != nil {
		return err
	}
	if quant.Quantity < purchase.Quantity {
		return fmt.Errorf("not enough products")
	}
	return repo.InsertQuantity(models.Quantity{
		ProductID: purchase.ProductID,
		Quantity:  quant.Quantity - purchase.Quantity,
	})
}

func productsPurchased(c *cli.Context) []dto.Purchase {
	result := make([]dto.Purchase, 0)

	for _, prod := range c.Args().Slice() {
		split := strings.Split(prod, "@")
		quantity := 1

		if len(split) > 1 {
			q, _ := strconv.Atoi(split[1])
			quantity += q
		}
		result = append(result, dto.Purchase{
			ProductID: split[0],
			Quantity:  quantity,
		})
	}
	return result
}

func validateBuyCommand(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("expected at least 1 product to buy")
	}
	pattern := regexp.MustCompile("[a-zA-Z0-9]+(@[0-9]+)?")
	for _, arg := range c.Args().Slice() {
		if !pattern.MatchString(arg) {
			return fmt.Errorf("arguments must be of the form: `productID` or `productID@count`")
		}
	}
	return nil
}

type notification struct {
	mode int
}

func (n *notification) IsPhone() bool {
	return n.mode == phoneNotification
}

func (n *notification) IsEmail()  bool {
	return n.mode == emailNotification
}
