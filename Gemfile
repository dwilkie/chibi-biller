source 'https://rubygems.org'
ruby "2.1.2"

# Bundle edge Rails instead: gem 'rails', github: 'rails/rails'
gem 'rails', '4.1.1'
gem 'resque'
gem 'unicorn'
gem 'httparty'

group :development, :test do
  gem 'foreman', :git => 'git://github.com/ddollar/foreman.git'
  gem 'rspec-rails'
end

group :test do
  gem 'webmock'
  gem 'vcr'
  gem 'resque_spec'
  gem 'capybara'
end
