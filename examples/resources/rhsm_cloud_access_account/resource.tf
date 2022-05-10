// Azure Example
resource "rhsm_cloud_access_account" "test_subscription" {
  account_id          = "123e4567-e89b-12d3-a456-426614174000"
  provider_short_name = "MSAZ"
  nickname            = "Test Azure Subscription"
  gold_images         = ["rhel-byos"]
}

// AWS Example
resource "rhsm_cloud_access_account" "test_account" {
  account_id          = "012345678912"
  provider_short_name = "AWS"
  nickname            = "Test AWS Account"
  gold_images         = ["RHEL"]
}

// GCP Example
resource "rhsm_cloud_access_account" "test_group" {
  account_id          = "test.group@example.com"
  provider_short_name = "GCE"
  nickname            = "Test GCP Group"
  gold_images         = ["rhel-byos-cloud-access"]
}
