import 'outfit.dart';

class OutfitRecommendation {
  final String id;
  final String outfitId;
  final Outfit? outfit;
  final double confidence;
  final String reasoning;
  final String? requestedOccasion;
  final String? requestedSeason;
  final String? requestedWeather;
  final String? requestedStyle;
  final bool viewed;
  final bool accepted;
  final DateTime? createdAt;

  OutfitRecommendation({
    required this.id,
    required this.outfitId,
    this.outfit,
    required this.confidence,
    required this.reasoning,
    this.requestedOccasion,
    this.requestedSeason,
    this.requestedWeather,
    this.requestedStyle,
    this.viewed = false,
    this.accepted = false,
    this.createdAt,
  });

  factory OutfitRecommendation.fromJson(Map<String, dynamic> json) {
    return OutfitRecommendation(
      id: json['id'] ?? '',
      outfitId: json['outfit_id'] ?? '',
      outfit: json['outfit'] != null 
          ? Outfit.fromJson(json['outfit']) 
          : null,
      confidence: json['confidence']?.toDouble() ?? 0.0,
      reasoning: json['reasoning'] ?? '',
      requestedOccasion: json['requested_occasion'],
      requestedSeason: json['requested_season'],
      requestedWeather: json['requested_weather'],
      requestedStyle: json['requested_style'],
      viewed: json['viewed'] ?? false,
      accepted: json['accepted'] ?? false,
      createdAt: json['created_at'] != null 
          ? DateTime.tryParse(json['created_at']) 
          : null,
    );
  }
}

class RecommendationRequest {
  final String? occasion;
  final String? season;
  final String? weather;
  final String? style;
  final List<String> colors;
  final int maxItems;

  RecommendationRequest({
    this.occasion,
    this.season,
    this.weather,
    this.style,
    this.colors = const [],
    this.maxItems = 5,
  });

  Map<String, dynamic> toJson() {
    return {
      'occasion': occasion,
      'season': season,
      'weather': weather,
      'style': style,
      'colors': colors,
      'max_items': maxItems,
    };
  }
}
