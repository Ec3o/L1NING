(function($) {
	$.fn.extend({
	  ProgressBarWars: function(options) {
		var $progressBar = this;
		var theidProgressBarWars = $progressBar.attr("id");
		var styleUnique = Date.now();
		var StringStyle = "";
		
		var defaults = {
		  porcentaje: "100",
		  tiempo: 1000,
		  color: "",
		  estilo: "yoda",
		  tamanio: "30%",
		  alto: "6px"
		};
		
		var settings = $.extend({}, defaults, options);
		if (settings.color !== '') {
		  StringStyle = "<style>.color" + styleUnique + 
			"{ border-radius: 2px; display: block; width: 0%; box-shadow: 0px 0px 10px 1px " + settings.color + 
			", 0 0 1px " + settings.color + ", 0 0 1px " + settings.color + 
			", 0 0 1px " + settings.color + ", 0 0 1px " + settings.color + 
			", 0 0 1px " + settings.color + ", 0 0 1px " + settings.color + 
			", 0 0 1px " + settings.color + "; background-color: #fff; }</style>";
		  settings.estilo = "color" + styleUnique;
		}
		
		$progressBar.before(StringStyle);
		$progressBar.append('<span class="barControl" style="width:' + settings.tamanio + 
		  ';"><div class="barContro_space"><span class="' + settings.estilo + 
		  '" style="height: ' + settings.alto + 
		  ';" id="bar' + theidProgressBarWars + '"></span></div></span>');
		$("#bar" + theidProgressBarWars).animate({ width: settings.porcentaje + "%" }, settings.tiempo);
		
		this.mover = function(ntamanio) {
		  $("#bar" + $progressBar.attr("id")).animate({ width: ntamanio + "%" }, settings.tiempo);
		};
		
		return this;
	  }
	});
  })(jQuery);
  