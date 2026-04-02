package jumpserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	Token         string
	BaseURL       string
	Username      string
	Password      string
	SkipTLSVerify bool
	APIVersion    string
}

func (c *Config) NewHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.SkipTLSVerify,
		},
	}
	return &http.Client{Transport: transport}
}

// GetAPIEndpoint returns the API v1 endpoint URL (used by most JumpServer resources)
func (c *Config) GetAPIEndpoint(path string) string {
	return fmt.Sprintf("%s/api/v1/%s", c.BaseURL, path)
}

// GetAPIEndpointV2 returns the API v2 endpoint URL (used by perms and other v2-only resources)
func (c *Config) GetAPIEndpointV2(path string) string {
	return fmt.Sprintf("%s/api/v2/%s", c.BaseURL, path)
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"skip_tls_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, skip SSL certificate validation (insecure).",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "v1",
				Deprecated:  "This field is deprecated and no longer affects API routing. Each resource now uses its correct API version automatically.",
				Description: "Deprecated. JumpServer API version (v1 or v2). No longer used — each resource targets its own API version.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jumpserver_user":             dataSourceUser(),
			"jumpserver_user_group":       dataSourceUserGroup(),
			"jumpserver_node":             dataSourceNode(),
			"jumpserver_zone":             dataSourceZone(),
			"jumpserver_platform":         dataSourcePlatform(),
			"jumpserver_asset":            dataSourceAsset(),
			"jumpserver_org":              dataSourceOrg(),
			"jumpserver_org_role":         dataSourceOrgRole(),
			"jumpserver_system_role":      dataSourceSystemRole(),
			"jumpserver_host":             dataSourceHost(),
			"jumpserver_label":            dataSourceLabel(),
			"jumpserver_account":          dataSourceAccount(),
			"jumpserver_account_template": dataSourceAccountTemplate(),
			"jumpserver_endpoint":         dataSourceEndpoint(),
			"jumpserver_database":         dataSourceDatabase(),
			"jumpserver_gateway":          dataSourceGateway(),
			"jumpserver_applet":           dataSourceApplet(),
			"jumpserver_device":           dataSourceDevice(),
			"jumpserver_web":              dataSourceWeb(),
			"jumpserver_cloud":            dataSourceCloud(),
			"jumpserver_custom":           dataSourceCustom(),
			"jumpserver_asset_permission": dataSourceAssetPermission(),
			"jumpserver_asset_category":   dataSourceAssetCategory(),
			"jumpserver_protocol":         dataSourceProtocol(),
			"jumpserver_content_type":     dataSourceContentType(),
			"jumpserver_terminal":         dataSourceTerminal(),
			"jumpserver_ticket_flow":      dataSourceTicketFlow(),
			"jumpserver_login_log":                dataSourceLoginLog(),
			"jumpserver_operate_log":              dataSourceOperateLog(),
			"jumpserver_activity_log":             dataSourceActivityLog(),
			"jumpserver_rbac_permission":           dataSourceRbacPermission(),
			"jumpserver_server_info":              dataSourceServerInfo(),
			"jumpserver_site_message":             dataSourceSiteMessage(),
			"jumpserver_sms_backend":              dataSourceSmsBackend(),
			"jumpserver_ops_task":                 dataSourceOpsTask(),
			"jumpserver_user_msg_subscription":    dataSourceUserMsgSubscription(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"jumpserver_host":             resourceHost(),
			"jumpserver_user":             resourceUser(),
			"jumpserver_user_group":       resourceUserGroup(),
			"jumpserver_node":             resourceNode(),
			"jumpserver_zone":             resourceZone(),
			"jumpserver_label":            resourceLabel(),
			"jumpserver_account":          resourceAccount(),
			"jumpserver_database":         resourceDatabase(),
			"jumpserver_device":           resourceDevice(),
			"jumpserver_web":              resourceWeb(),
			"jumpserver_gateway":          resourceGateway(),
			"jumpserver_cloud":            resourceCloud(),
			"jumpserver_custom":           resourceCustom(),
			"jumpserver_asset":              resourceAsset(),
			"jumpserver_system_user":        resourceSystemUser(),
			"jumpserver_asset_permission":   resourceAssetPermission(),
			"jumpserver_command_group":      resourceCommandGroup(),
			"jumpserver_command_filter_acl": resourceCommandFilterACL(),
			"jumpserver_login_acl":          resourceLoginACL(),
			"jumpserver_login_asset_acl":    resourceLoginAssetACL(),
			"jumpserver_connect_method_acl": resourceConnectMethodACL(),
			"jumpserver_data_masking_rule":          resourceDataMaskingRule(),
			"jumpserver_account_template":           resourceAccountTemplate(),
			"jumpserver_account_backup_plan":        resourceAccountBackupPlan(),
			"jumpserver_integration_application":    resourceIntegrationApplication(),
			"jumpserver_endpoint":                   resourceEndpoint(),
			"jumpserver_endpoint_rule":              resourceEndpointRule(),
			"jumpserver_command_storage":            resourceCommandStorage(),
			"jumpserver_replay_storage":             resourceReplayStorage(),
			"jumpserver_org_role":                   resourceOrgRole(),
			"jumpserver_system_role":                resourceSystemRole(),
			"jumpserver_role_binding":               resourceRoleBinding(),
			"jumpserver_job":                        resourceJob(),
			"jumpserver_playbook":                   resourcePlaybook(),
			"jumpserver_directory":                  resourceDirectory(),
			"jumpserver_gpt":                        resourceGPT(),
			"jumpserver_platform":                   resourcePlatform(),
			"jumpserver_org":                        resourceOrg(),
			"jumpserver_org_role_binding":            resourceOrgRoleBinding(),
			"jumpserver_system_role_binding":         resourceSystemRoleBinding(),
			"jumpserver_ticket_flow":                resourceTicketFlow(),
			"jumpserver_adhoc":                      resourceAdHoc(),
			"jumpserver_ops_variable":               resourceOpsVariable(),
			"jumpserver_access_key":                 resourceAccessKey(),
			"jumpserver_virtual_account":            resourceVirtualAccount(),
			"jumpserver_labeled_resource":            resourceLabeledResource(),
			"jumpserver_applet":                     resourceApplet(),
			"jumpserver_applet_host":                resourceAppletHost(),
			"jumpserver_applet_publication":          resourceAppletPublication(),
			"jumpserver_virtual_app":                resourceVirtualApp(),
			"jumpserver_virtual_app_publication":     resourceVirtualAppPublication(),
			"jumpserver_ssh_key":                     resourceSSHKey(),
			"jumpserver_user_group_relation":         resourceUserGroupRelation(),
			"jumpserver_passkey":                     resourcePasskey(),
			"jumpserver_protocol_setting":            resourceProtocolSetting(),
			"jumpserver_leak_password":               resourceLeakPassword(),
			"jumpserver_app_provider":                resourceAppProvider(),
			"jumpserver_applet_host_deployment":      resourceAppletHostDeployment(),
			"jumpserver_account_risk":                resourceAccountRisk(),
			"jumpserver_gathered_account":            resourceGatheredAccount(),
			"jumpserver_favorite_asset":              resourceFavoriteAsset(),
			"jumpserver_asset_perm_user_relation":    resourceAssetPermUserRelation(),
			"jumpserver_asset_perm_user_group_relation": resourceAssetPermUserGroupRelation(),
			"jumpserver_asset_perm_asset_relation":   resourceAssetPermAssetRelation(),
			"jumpserver_asset_perm_node_relation":    resourceAssetPermNodeRelation(),
			"jumpserver_ticket":                      resourceTicket(),
			"jumpserver_apply_asset_ticket":          resourceApplyAssetTicket(),
			"jumpserver_apply_command_ticket":        resourceApplyCommandTicket(),
			"jumpserver_apply_login_ticket":          resourceApplyLoginTicket(),
			"jumpserver_apply_login_asset_ticket":    resourceApplyLoginAssetTicket(),
			"jumpserver_ticket_comment":              resourceTicketComment(),
			"jumpserver_terminal":                    resourceTerminal(),
			"jumpserver_session":                     resourceSession(),
			"jumpserver_session_sharing":             resourceSessionSharing(),
			"jumpserver_session_join_record":         resourceSessionJoinRecord(),
			"jumpserver_terminal_command":            resourceTerminalCommand(),
			"jumpserver_chatai_prompt":               resourceChatAIPrompt(),
			"jumpserver_connection_token":            resourceConnectionToken(),
			"jumpserver_rbac_role":                    resourceRbacRole(),
			"jumpserver_setting":                      resourceSetting(),
			"jumpserver_super_connection_token":       resourceSuperConnectionToken(),
			"jumpserver_temp_token":                   resourceTempToken(),
			"jumpserver_terminal_registration":        resourceTerminalRegistration(),
			"jumpserver_system_msg_subscription":      resourceSystemMsgSubscription(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	skipTLS := d.Get("skip_tls_verify").(bool)
	apiVersion := d.Get("api_version").(string)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLS,
		},
	}
	client := &http.Client{Transport: transport}

	token, err := getToken(client, baseURL, username, password, apiVersion)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return &Config{
		Token:         token,
		BaseURL:       baseURL,
		Username:      username,
		Password:      password,
		SkipTLSVerify: skipTLS,
		APIVersion:    apiVersion,
	}, diags
}

func getToken(client *http.Client, baseURL, username, password, apiVersion string) (string, error) {
	// Always authenticate via /api/v1/ (works on all JumpServer versions)
	url := baseURL + "/api/v1/authentication/auth/"

	credentials := map[string]string{
		"username": username,
		"password": password,
	}
	jsonValue, _ := json.Marshal(credentials)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if token, ok := result["token"].(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("unable to fetch token from %s", url)
}
