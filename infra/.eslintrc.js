module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  env: {
    node: true
  },
  plugins: ["@typescript-eslint", "prettier"],
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/eslint-recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier/@typescript-eslint"
  ],
  rules: {
    "prettier/prettier": "error"
  }
};
