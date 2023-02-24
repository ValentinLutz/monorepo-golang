const openAPIFiles = [
  {
    url: '../openapi/order_api.yaml',
    name: 'Order API'
  }
];

window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    urls: openAPIFiles,
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    displayRequestDuration: true,
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
