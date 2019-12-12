package sealedsecrets

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rockyhmchen/terraform-provider-sealedsecrets/utils"
)

type Cmd struct {
	kubectl  string
	kubeseal string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kubectl_bin": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     utils.Which("kubectl"),
				Description: "TBA",
			},

			"kubeseal_bin": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     utils.Which("kubeseal"),
				Description: "TBA",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sealedsecrets_secret": resourceSecret(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	kubectl := d.Get("kubectl_bin").(string)
	kubeseal := d.Get("kubeseal_bin").(string)

	if !utils.PathExists(kubectl) {
		errMsg := "kubectl command doesn't exist"
		utils.Log(errMsg)
		return nil, errors.New(errMsg)
	}

	if !utils.PathExists(kubeseal) {
		errMsg := "kubeseal command doesn't exist"
		utils.Log(errMsg)
		return nil, errors.New(errMsg)
	}

	return &Cmd{kubectl, kubeseal}, nil
}
