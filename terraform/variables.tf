variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "Region for resources"
  type        = string
  default     = "us-central1"
}

variable "lab_id" {
  description = "Unique identifier for the lab instance"
  type        = string
}

variable "ssh_public_key_path" {
  description = "Path to SSH public key"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

variable "dummy_mode" {
  description = "Use dummy provider instead of real infrastructure"
  type        = bool
  default     = true
}

variable "debug" {
  description = "Enable debug mode with dummy resources"
  type        = bool
  default     = false
}
