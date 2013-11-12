source 'https://rubygems.org'
ruby "2.0.0"

# Bundle edge Rails instead: gem 'rails', github: 'rails/rails'
gem 'rails', '4.0.1'
gem 'resque'
gem 'unicorn'
gem 'httparty'

group :development, :test do
  gem 'foreman', :git => 'git://github.com/ddollar/foreman.git'
  gem 'rspec-rails'
end

group :test do
  gem 'spork-rails', :git => 'git://github.com/sporkrb/spork-rails.git'
  gem 'webmock'
  gem 'vcr', :git => 'git://github.com/myronmarston/vcr.git'
  gem 'guard-rspec'
  gem 'guard-spork'
  gem 'rb-inotify'
  gem 'resque_spec'
  gem 'capybara'
end
