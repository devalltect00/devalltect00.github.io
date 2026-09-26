import eslintPluginAstro from "eslint-plugin-astro";
import tseslint from "typescript-eslint";

export default [
  ...tseslint.configs.recommended,
  ...eslintPluginAstro.configs["flat/recommended"],
  {
    ignores: [
      ".astro/**",
      "dist/**",
      "node_modules/**",
      "docs/project_structure.md",
    ],
  },
];
