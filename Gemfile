source 'https://rubygems.org'

# Bundle edge Rails instead: gem 'rails', github: 'rails/rails'
gem 'rails', '4.0.0'
gem 'resque'
gem "unicorn"

group :development, :test do
  gem 'foreman'
  gem 'rspec-rails'
end

group :test do
  gem 'sqlite3'
  gem 'spork-rails', :git => 'git://github.com/sporkrb/spork-rails.git'
  gem 'mock_redis'
  gem 'webmock'
  gem 'vcr', :git => 'git://github.com/myronmarston/vcr.git'
  gem 'guard-rspec'
  gem 'guard-spork'
  gem 'rb-inotify'
  gem 'resque_spec'
  gem 'capybara'
end
