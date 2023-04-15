const express = require("express");
const swaggerUi = require("swagger-ui-express");
const YAML = require("yamljs");
const swaggerDocument = YAML.load("./openapi.yaml");

const app = express();
const PORT = 3001;

app.use("/docs", swaggerUi.serve, swaggerUi.setup(swaggerDocument));

app.get(
  "/",
  swaggerUi.setup(null, {
    swaggerOptions: {
      requestInterceptor: function (request) {
        request.headers.Origin = `http://localhost:3001`;
        return request;
      },
      url: `http://localhost:3001/docs/`,
    },
  })
);

app.listen(PORT, () => {
  console.log(`Swagger UI is available at http://localhost:${PORT}`);
});
