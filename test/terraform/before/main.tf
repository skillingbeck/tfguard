resource "time_rotating" "t1" {
  rotation_days = 30
  triggers = {
      version = 1
  }
}

resource "time_rotating" "t2" {
  rotation_days = 30
  triggers = {
      version = 1
  }
}

resource "time_rotating" "t3" {
  rotation_days = 30
  triggers = {
      version = 1
  }
}