<template>
    <div>
      <div ref="pieChart" :style="{ width: '100%', height: '400px' }"></div>
      <div class="data-summary">
        <p class="cool-font">total mileage: {{ chartData.total_mileage }}</p>
        <p class="cool-font">valid mileage: {{ chartData.valid_mileage }}</p>
        <p class="cool-font">average speed: {{ chartData.average_speed }}</p>
        <p class="cool-font">valid times: {{ chartData.valid_times }}</p>
        <p class="cool-font">Total: 24</p>
      </div>
    </div>
  </template>
  
  <script>
  import { nextTick } from 'vue';
  import * as echarts from 'echarts';
  
  export default {
    name: 'EchartsPieChart',
    props: {
      chartData: {
        type: Object,
        required: true
      }
    },
    watch: {
      chartData: {
        handler() {
          this.renderChart();
        },
        deep: true,
        immediate: true
      }
    },
    mounted() {
      this.renderChart();
    },
    methods: {
      renderChart() {
        nextTick(() => {
          if (!this.$refs.pieChart) {
            return;
          }
          const chart = echarts.init(this.$refs.pieChart);
          const validTimes = parseInt(this.chartData.valid_times, 10);
          const invalidTimes = 24 - validTimes;
  
          const option = {
            tooltip: {
              trigger: 'item'
            },
            legend: {
              orient: 'horizontal',
              bottom: '0%',
              textStyle: {
                fontSize: 14
              }
            },
            series: [
              {
                name: 'Status',
                type: 'pie',
                radius: '50%',
                data: [
                  { value: validTimes, name: 'Valid' },
                  { value: invalidTimes, name: 'Remain' }
                ],
                emphasis: {
                  itemStyle: {
                    shadowBlur: 10,
                    shadowOffsetX: 0,
                    shadowColor: 'rgba(0, 0, 0, 0.5)'
                  }
                },
                label: {
                  formatter: '{b|{b}}\n{d|{d}%}',
                  rich: {
                    b: {
                      fontSize: 14,
                      lineHeight: 20
                    },
                    d: {
                      fontSize: 14,
                      lineHeight: 20,
                      align: 'center'
                    }
                  }
                }
              }
            ]
          };
          chart.setOption(option);
        });
      }
    }
  };
  </script>
  
  <style scoped>
  .cool-font {
    font-family: 'Pacifico', cursive;
  }
  .data-summary {
    margin-top: 20px;
    font-size: 16px;
  }
  </style>
  