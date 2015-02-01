source 'https://rubygems.org'
ruby "2.2.0"

# Bundle edge Rails instead: gem 'rails', github: 'rails/rails'
gem 'rails', '4.2.0'
gem 'unicorn'
gem 'httparty'
gem 'sidekiq'

group :development do
  gem 'capistrano'
  gem 'capistrano-rails'
  gem 'capistrano-rbenv'
  gem 'capistrano-bundler'
  gem 'capistrano3-foreman'
end

group :deployment, :development do
  gem 'foreman', :github => "ddollar/foreman" # needed for deployment
end

group :development, :test do
  gem 'rspec-rails'
  gem 'pry'
end

group :test do
  gem 'webmock'
  gem 'vcr'
  gem 'capybara'
end
