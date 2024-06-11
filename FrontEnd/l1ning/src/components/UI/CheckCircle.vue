<template>
  <svg id="svg" width="400" height="400" @click="toggleCheck">
    <circle
      id="circle"
      fill="grey"
      stroke="black"
      stroke-width="20"
      cx="200"
      cy="200"
      r="190"
      stroke-linecap="round"
      transform="rotate(-90 200 200)"
    />
    <polyline
      id="tick-border"
      fill="none"
      stroke="black"
      stroke-width="32"
      points="88,214 173,284 304,138"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <polyline
      id="tick"
      fill="none"
      stroke="white"
      stroke-width="24"
      points="88,214 173,284 304,138"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>

<script>
export default {
  name: "CheckCircle",
  props: {
    checked: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      size: 40
    };
  },
  watch: {
    checked(newVal) {
      this.updateCheckState(newVal);
    }
  },
  mounted() {
    this.initialize();
    this.updateCheckState(this.checked);
  },
  methods: {
    initialize() {
      const size = this.size;
      const svg = this.$el;
      const circle = svg.querySelector("#circle");
      const tick = svg.querySelector("#tick");
      const tickBorder = svg.querySelector("#tick-border");

      svg.setAttribute("width", size);
      svg.setAttribute("height", size);

      const strokeWidth = parseFloat(size) / 20;
      circle.setAttribute("stroke-width", strokeWidth);
      circle.setAttribute("cx", parseFloat(size) / 2);
      circle.setAttribute("cy", parseFloat(size) / 2);
      circle.setAttribute("r", parseFloat(size) / 2 - strokeWidth / 2);
      circle.setAttribute(
        "transform",
        `rotate(-90 ${parseFloat(size) / 2} ${parseFloat(size) / 2})`
      );
      circle.setAttribute("stroke-dasharray", parseFloat(size) / 4);
      circle.setAttribute("stroke-dashoffset", 11.94 * parseFloat(size) / 4);

      const tickStrokeWidth = parseFloat(size) / 16.67;
      tick.setAttribute("stroke-width", tickStrokeWidth);
      tickBorder.setAttribute("stroke-width", tickStrokeWidth + 4);
      const points = `${parseFloat(size) * 0.22},${parseFloat(size) * 0.54} ${parseFloat(size) * 0.43},${parseFloat(size) * 0.71} ${parseFloat(size) * 0.76},${parseFloat(size) * 0.34}`;
      tick.setAttribute("points", points);
      tickBorder.setAttribute("points", points);
    },
    updateCheckState(checked) {
      const size = this.size;
      const svg = this.$el;
      const circle = svg.querySelector("#circle");
      const tick = svg.querySelector("#tick");

      if (checked) {
        circle.setAttribute("stroke-dasharray", 3 * parseFloat(size));
        circle.setAttribute("stroke-dashoffset", 0);
        circle.setAttribute("fill", "green");
        tick.setAttribute("stroke", "white");
        svg.classList.add("checked");
      } else {
        circle.setAttribute("stroke-dasharray", parseFloat(size) / 4);
        circle.setAttribute("stroke-dashoffset", 11.94 * parseFloat(size) / 4);
        circle.setAttribute("fill", "grey");
        tick.setAttribute("stroke", "white");
        svg.classList.remove("checked");
      }
    },
    toggleCheck() {
      const newCheckedState = !this.checked;
      this.$emit('update:checked', newCheckedState);
      this.updateCheckState(newCheckedState);
    }
  }
};
</script>

<style scoped>
:root {
  --size: 40px;
}

body {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  flex-direction: column;
}

#circle {
  stroke-dasharray: 10;
  stroke-dashoffset: 11.94;
  transition: stroke-dashoffset 1s ease-in-out, fill 0.5s ease-in-out;
}

#tick,
#tick-border {
  stroke-dasharray: 350;
  stroke-dashoffset: 350;
  transition: stroke-dashoffset 1.6s ease-out;
}

#svg {
  transition: transform 1s ease-in-out;
}

.checked #circle {
  stroke-dasharray: calc(3 * var(--size));
  stroke-dashoffset: 0;
  fill: green;
}

.checked #tick,
.checked #tick-border {
  stroke-dashoffset: 0;
}
</style>
