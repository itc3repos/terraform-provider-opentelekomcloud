package opentelekomcloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/auto_recovery"
)

func resourceECSAutoRecoveryV1Read(d *schema.ResourceData, meta interface{}, instanceID string) (bool, error) {
	config := meta.(*Config)
	client, err := chooseECSV1Client(d, config)
	if err != nil {
		return false, fmt.Errorf("Error creating OpenTelekomCloud client: %s", err)
	}

	rId := instanceID

	r, err := auto_recovery.Get(client, rId).Extract()
	if err != nil {
		return false, CheckDeleted(d, err, "ECS-AutoRecovery")
	}
	log.Printf("[DEBUG] Retrieved ECS-AutoRecovery:%#v of instance:%s", rId, r)
	return strconv.ParseBool(r.SupportAutoRecovery)
}

func setAutoRecoveryForInstance(d *schema.ResourceData, meta interface{}, instanceID string, ar bool) error {
	config := meta.(*Config)
	client, err := chooseECSV1Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud client: %s", err)
	}

	rId := instanceID

	updateOpts := auto_recovery.UpdateOpts{SupportAutoRecovery: strconv.FormatBool(ar)}

	timeout := d.Timeout(schema.TimeoutUpdate)

	log.Printf("[DEBUG] Setting ECS-AutoRecovery for instance:%s with options: %#v", rId, updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := auto_recovery.Update(client, rId, updateOpts)
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error setting ECS-AutoRecovery for instance%s: %s", rId, err)
	}
	return nil
}
