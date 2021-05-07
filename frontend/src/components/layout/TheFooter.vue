<template>
  <footer>
    <nav>
      <h1>
        <img src="@/assets/logo.png" alt="mrratality" />
        <router-link to="/">mrratality</router-link>
      </h1>
      <ul>
        <li>© 2021 hackfeed</li>
        <li><a href="mailto::hackfeed@yandex.ru">Contact me</a></li>
        <li><a href="https://github.com/hackfeed/mrratality">Project Page</a></li>
      </ul>
    </nav>
  </footer>
</template>

<script>
export default {};
</script>

<style scoped>
footer {
  width: 100%;
  height: 3rem;
  border-top: 1px solid rgba(0, 0, 0, 0.26);
  position: absolute;
  bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
}
li a {
  color: black;
  text-decoration: none;
}
li a:hover {
  text-decoration: underline;
}
h1 {
  margin: 0;
  display: flex;
  align-items: center;
  color: black;
}
h1 img {
  height: 2rem;
}
h1 a {
  color: black;
  text-decoration: none;
  margin: 0;
  padding: 0.75rem 1.5rem;
  display: inline-block;
}
h1 a:hover,
h1 a:active,
h1 a.router-link-active {
  border-color: transparent;
}
footer nav {
  width: 90%;
  margin: auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
footer ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  justify-content: center;
  align-items: center;
}
li {
  margin: 0 0.5rem;
}
</style>
