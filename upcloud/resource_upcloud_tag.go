package upcloud

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUpCloudTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceUpCloudTagCreate,
		Update: resourceUpCloudTagUpdate,
		Delete: resourceUpCloudTagDelete,
		Read:   resourceUpCloudTagRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceUpCloudTagRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUpCloudTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)

	createTagRequest := &request.CreateTagRequest{
		Tag: upcloud.Tag{
			Name: d.Get("name").(string),
		},
	}
	if description, ok := d.GetOk("description"); ok {
		createTagRequest.Description = description.(string)
	}
	if instances, ok := d.GetOk("instances"); ok {
		instances := instances.([]interface{})
		instancesList := make([]string, len(instances))
		for i := range instancesList {
			instancesList[i] = instances[i].(string)
		}
		createTagRequest.Servers = instancesList
	}

	tag, err := client.CreateTag(createTagRequest)

	if err != nil {
		return err
	}

	d.SetId(tag.Name)

	return nil
}

func resourceUpCloudTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)
	r := &request.ModifyTagRequest{
		Name: d.Id(),
	}

	r.Tag.Name = d.Id()
	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		r.Tag.Description = newDescription.(string)
	}
	if d.HasChange("instances") {
		_, newServers := d.GetChange("instances")

		instances := newServers.([]interface{})
		instancesList := make([]string, len(instances))
		for i := range instancesList {
			instancesList[i] = instances[i].(string)
		}
		r.Tag.Servers = instancesList
	}

	_, err := client.ModifyTag(r)

	if err != nil {
		return err
	}

	return nil
}

func resourceUpCloudTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)

	deleteTagRequest := &request.DeleteTagRequest{
		Name: d.Id(),
	}
	err := client.DeleteTag(deleteTagRequest)

	if err != nil {
		return err
	}

	return nil
}
