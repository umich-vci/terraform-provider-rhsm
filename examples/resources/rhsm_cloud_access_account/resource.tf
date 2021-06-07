resource "rhsm_cloud_access_account" "test_account" {
  account_id          = "123e4567-e89b-12d3-a456-426614174000"
  provider_short_name = "MSAZ"
  nickname            = "Test Account"
  gold_images         = ["rhel-byos"]
}
