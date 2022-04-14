import { Selector } from "testcafe"

fixture`Getting started`
  .page`http://localhost:8081`

test("first test", () => {
})


test("has jquery", async t => {
  await t
    .expect(Selector("#jquery").exists).ok()
})
