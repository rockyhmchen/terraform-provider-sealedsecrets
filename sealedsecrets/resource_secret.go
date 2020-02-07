package sealedsecrets

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/rockyhmchen/terraform-provider-sealedsecrets/utils"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "TBA",
			},
			"namespace": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "TBA",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "TBA",
			},
			"secret_source": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "TBA",
			},
			"sealed_secret_source": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "TBA",
			},
			"certificate": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "TBA",
			},
		},
	}
}

func resourceSecretCreate(d *schema.ResourceData, m interface{}) error {
	utils.Log("resourceSecretCreate")

	mainCmd := m.(*Cmd)

	ssPath := d.Get("sealed_secret_source").(string)
	if shouldCreateSealedSecret(ssPath) {
		utils.Log("Sealed secret doesn't exist")
		if err := createSealedSecret(d, mainCmd); err != nil {
			return err
		}
	}
	utils.Log(fmt.Sprintf("Sealed secret (%s) has been created\n", ssPath))

	if err := utils.ExecuteCmd(mainCmd.kubectl, "apply", "-f", ssPath); err != nil {
		return err
	}

	d.SetId(utils.SHA256(utils.GetFileContent(ssPath)))

	return resourceSecretRead(d, m)
}

func resourceSecretDelete(d *schema.ResourceData, m interface{}) error {
	utils.Log("resourceSecretDelete")

	mainCmd := m.(*Cmd)

	name := d.Get("name").(string)
	ns := d.Get("namespace").(string)

	if err := utils.ExecuteCmd(mainCmd.kubectl, "delete", "secret", name, "-n", ns); err != nil {
		return err
	}

	if err := utils.ExecuteCmd(mainCmd.kubectl, "delete", "SealedSecret", name, "-n", ns); err != nil {
		return err
	}

	// delete sealed secrets file
	ssPath := d.Get("sealed_secret_source").(string)
	if err := os.Remove(ssPath); err != nil {
		utils.Log("Failed to delete sealed secret file: " + ssPath)
	}

	d.SetId("")

	return nil
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	utils.Log("resourceSecretRead")

	ssPath := d.Get("sealed_secret_source").(string)
	if shouldCreateSealedSecret(ssPath) {
		d.SetId("")
		return nil
	}

	oldID := d.Id()
	newID := utils.SHA256(utils.GetFileContent(ssPath))
	if oldID != newID {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, m interface{}) error {
	utils.Log("resourceSecretUpdate")

	mainCmd := m.(*Cmd)

	name := d.Get("name").(string)
	ns := d.Get("namespace").(string)

	if err := utils.ExecuteCmd(mainCmd.kubectl, "delete", "secret", name, "-n", ns); err != nil {
		utils.Log(fmt.Sprintf("Failed to delete secret: %s.%s\n", ns, name))
	}

	if err := utils.ExecuteCmd(mainCmd.kubectl, "delete", "SealedSecret", name, "-n", ns); err != nil {
		utils.Log(fmt.Sprintf("Failed to delete SealedSecret: %s.%s\n", ns, name))
	}

	// delete sealed secrets file
	ssPath := d.Get("sealed_secret_source").(string)
	if err := os.Remove(ssPath); err != nil {
		utils.Log("Failed to delete sealed secret file: " + ssPath)
	}

	return resourceSecretCreate(d, m)
}

func shouldCreateSealedSecret(ssPath string) bool {
	return !utils.PathExists(ssPath)
}

func createSealedSecret(d *schema.ResourceData, mainCmd *Cmd) error {
	sPath := d.Get("secret_source").(string)
	if !utils.PathExists(sPath) {
		errMsg := "Could not find secret source"
		utils.Log(errMsg)
		return errors.New(errMsg)
	}

	cert := d.Get("certificate").(string)
	if !utils.PathExists(cert) {
		errMsg := "Could not find certificate"
		utils.Log(errMsg)
		return errors.New(errMsg)
	}

	ssPath := d.Get("sealed_secret_source").(string)
	ssDir := utils.GetDir(ssPath)
	if !utils.PathExists(ssDir) {
		utils.Log(fmt.Sprintf("Sealed secret directory (%s) doesn't exist\n", ssDir))
		os.Mkdir(ssDir, os.ModePerm)
	}
	utils.Log(fmt.Sprintf("Sealed secret (%s) has been created\n", ssDir))

	name := d.Get("name").(string)
	ns := d.Get("namespace").(string)
	sType := d.Get("type").(string)

	nsArg := fmt.Sprintf("%s=%s", "--namespace", ns)
	typeArg := fmt.Sprintf("%s=%s", "--type", sType)
	fromFileArg := fmt.Sprintf("%s=%s=%s", "--from-file", utils.GetFileName(sPath), sPath)
	dryRunArg := "--dry-run"
	outputArg := fmt.Sprintf("%s=%s", "--output", "yaml")

	certArg := fmt.Sprintf("%s %s", "--cert", cert)
	formatArg := fmt.Sprintf("%s %s", "--format", "yaml")

	return utils.ExecuteCmd(mainCmd.kubectl, "create", "secret", "generic", name, nsArg, typeArg, fromFileArg, dryRunArg, outputArg, " | ", mainCmd.kubeseal, ">", ssPath, certArg, formatArg)
}
