
# Service apiv3\n
				
# Service apiv3-cache\n
				
# Service apiv3-outbound\n
				
# Service apiv3-rack\n
				
# Service apiv3-redis\n
				
# Service apiv3-sequel\n
				
import {
  id = "apiv3-sequel"
  to = module.apiv3_sequel.datadog_service_definition_yaml.service_definition[0]
}

					 
# Service apiv3-sidekiq\n
				
# Service apiv3-sinatra\n
				
import {
  id = "apiv3-sinatra"
  to = module.apiv3_sinatra.datadog_service_definition_yaml.service_definition[0]
}

					 
# Service classic\n
				
# Service classic-activerecord\n
				
# Service classic-redis\n
				
# Service classic-view\n
				
# Service classic-controller\n
				
# Service classic-outbound\n
				
# Service classic-resque\n
				
# Service authentication-api\n
				
import {
  id = "teamsnap/authentication-api"
  to = module.authentication_api.module.sentry[0].sentry_project.this
}

					 
# Service authentication-service\n
				
import {
  id = "teamsnap/authentication-service"
  to = module.authentication_service.module.sentry[0].sentry_project.this
}

					 
# Service person-service\n
				
import {
  id = "teamsnap/person-service"
  to = module.person_service.module.sentry[0].sentry_project.this
}

					 
# Service user-account-service\n
				
import {
  id = "teamsnap/user-account-service"
  to = module.user_account_service.module.sentry[0].sentry_project.this
}

					 
# Service cogsworth\n
				
# Service snapi\n
				