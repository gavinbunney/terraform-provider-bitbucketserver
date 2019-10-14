/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

// See https://docusaurus.io/docs/site-config for all the possible
// site configuration options.

const siteConfig = {
  title: 'Bitbucket Server Terraform Provider',

  url: 'https://gavinbunney.github.io',
  baseUrl: '/',
  // baseUrl: '/terraform-provider-bitbucketserver/',

  // Used for publishing and more
  projectName: 'terraform-provider-bitbucketserver',
  organizationName: 'gavinbunney',

  // For no header links in the top nav bar -> headerLinks: [],
  headerLinks: [
  ],

  /* path to images for header/footer */
  favicon: 'img/favicon.png',

  /* Colors for website */
  colors: {
    primaryColor: '#5C4EE5',
    secondaryColor: '#66246c',
  },

  // This copyright info is used in /core/Footer.js and blog RSS/Atom feeds.
  copyright: `Copyright © ${new Date().getFullYear()} Gavin Bunney`,

  highlight: {
    // Highlight.js theme to use for syntax highlighting in code blocks.
    theme: 'default',
  },

  // Add custom scripts here that would be placed in <script> tags.
  scripts: ['https://buttons.github.io/buttons.js'],

  stylesheets: [
    "https://fonts.googleapis.com/css?family=Lato:400,400i,700,700i,900,900i&display=swap"
  ],

  // On page navigation for the current documentation page.
  onPageNav: 'separate',

  // No .html extensions for paths.
  cleanUrl: true,
};

module.exports = siteConfig;
