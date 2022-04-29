<script lang='ts'>
</script>

<nav>
  <slot></slot>
</nav>

<style>
  nav {
    position: relative;
    width: 100%;
    display: grid;
    grid-template-rows: minmax(0, 1fr);
    justify-content: flex-start;
    border-bottom: 1px solid teal;
  }
</style>