<script lang="ts">
  export let title: string;
  export let desc: string;
  export let type: "warn";
</script>

<div class={type}>
  <h4 class="norm">{title}</h4>
  <p>{desc}</p>
</div>

<style lang="scss">
  div {
    padding: 10px 12px;
    border-radius: 10px;

    &.warn {
      background-color: $warn;
    }
  }
</style>
