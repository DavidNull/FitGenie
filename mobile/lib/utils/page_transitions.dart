import 'package:flutter/material.dart';

/// Smooth fade transition between pages
class FadePageRoute<T> extends PageRouteBuilder<T> {
  final Widget child;
  final Duration duration;

  FadePageRoute({
    required this.child,
    this.duration = const Duration(milliseconds: 300),
  }) : super(
          pageBuilder: (context, animation, secondaryAnimation) => child,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            return FadeTransition(
              opacity: animation,
              child: child,
            );
          },
          transitionDuration: duration,
        );
}

/// Slide transition from right (iOS style)
class SlideRightPageRoute<T> extends PageRouteBuilder<T> {
  final Widget child;
  final Duration duration;

  SlideRightPageRoute({
    required this.child,
    this.duration = const Duration(milliseconds: 300),
  }) : super(
          pageBuilder: (context, animation, secondaryAnimation) => child,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            const begin = Offset(1.0, 0.0);
            const end = Offset.zero;
            const curve = Curves.easeInOutCubic;

            var tween = Tween(begin: begin, end: end).chain(
              CurveTween(curve: curve),
            );

            return SlideTransition(
              position: animation.drive(tween),
              child: FadeTransition(
                opacity: animation,
                child: child,
              ),
            );
          },
          transitionDuration: duration,
        );
}

/// Slide transition from bottom (modal style)
class SlideUpPageRoute<T> extends PageRouteBuilder<T> {
  final Widget child;
  final Duration duration;

  SlideUpPageRoute({
    required this.child,
    this.duration = const Duration(milliseconds: 400),
  }) : super(
          pageBuilder: (context, animation, secondaryAnimation) => child,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            const begin = Offset(0.0, 1.0);
            const end = Offset.zero;
            const curve = Curves.easeOutCubic;

            var tween = Tween(begin: begin, end: end).chain(
              CurveTween(curve: curve),
            );

            return SlideTransition(
              position: animation.drive(tween),
              child: child,
            );
          },
          transitionDuration: duration,
        );
}

/// Scale transition with fade
class ScalePageRoute<T> extends PageRouteBuilder<T> {
  final Widget child;
  final Duration duration;

  ScalePageRoute({
    required this.child,
    this.duration = const Duration(milliseconds: 300),
  }) : super(
          pageBuilder: (context, animation, secondaryAnimation) => child,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            const curve = Curves.easeOutCubic;

            var scaleTween = Tween(begin: 0.9, end: 1.0).chain(
              CurveTween(curve: curve),
            );

            return FadeTransition(
              opacity: animation,
              child: ScaleTransition(
                scale: animation.drive(scaleTween),
                child: child,
              ),
            );
          },
          transitionDuration: duration,
        );
}

/// Helper class for common transitions
class PageTransitions {
  /// Navigate with fade transition
  static Future<T?> fadeTo<T>(BuildContext context, Widget page) {
    return Navigator.push<T>(
      context,
      FadePageRoute(child: page),
    );
  }

  /// Navigate with slide from right (iOS style)
  static Future<T?> slideTo<T>(BuildContext context, Widget page) {
    return Navigator.push<T>(
      context,
      SlideRightPageRoute(child: page),
    );
  }

  /// Navigate with slide from bottom (modal)
  static Future<T?> modalTo<T>(BuildContext context, Widget page) {
    return Navigator.push<T>(
      context,
      SlideUpPageRoute(child: page),
    );
  }

  /// Navigate with scale + fade
  static Future<T?> scaleTo<T>(BuildContext context, Widget page) {
    return Navigator.push<T>(
      context,
      ScalePageRoute(child: page),
    );
  }
}

/// Animated route for hero-like transitions
class SharedAxisPageRoute<T> extends PageRouteBuilder<T> {
  final Widget child;
  final SharedAxisTransitionType type;
  final Duration duration;

  SharedAxisPageRoute({
    required this.child,
    this.type = SharedAxisTransitionType.scaled,
    this.duration = const Duration(milliseconds: 400),
  }) : super(
          pageBuilder: (context, animation, secondaryAnimation) => child,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            return SharedAxisTransition(
              animation: animation,
              secondaryAnimation: secondaryAnimation,
              transitionType: type,
              child: child,
            );
          },
          transitionDuration: duration,
        );
}

/// Material shared axis transition (simplified)
class SharedAxisTransition extends StatelessWidget {
  final Animation<double> animation;
  final Animation<double> secondaryAnimation;
  final SharedAxisTransitionType transitionType;
  final Widget child;

  const SharedAxisTransition({
    super.key,
    required this.animation,
    required this.secondaryAnimation,
    required this.transitionType,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    final fadeAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).chain(CurveTween(curve: const Interval(0.3, 1.0))).animate(animation);

    switch (transitionType) {
      case SharedAxisTransitionType.scaled:
        final scaleAnimation = Tween<double>(
          begin: 0.85,
          end: 1.0,
        ).chain(CurveTween(curve: Curves.easeOutCubic)).animate(animation);

        return FadeTransition(
          opacity: fadeAnimation,
          child: ScaleTransition(
            scale: scaleAnimation,
            child: child,
          ),
        );

      case SharedAxisTransitionType.horizontal:
        final slideAnimation = Tween<Offset>(
          begin: const Offset(0.1, 0.0),
          end: Offset.zero,
        ).chain(CurveTween(curve: Curves.easeOutCubic)).animate(animation);

        return FadeTransition(
          opacity: fadeAnimation,
          child: SlideTransition(
            position: slideAnimation,
            child: child,
          ),
        );

      case SharedAxisTransitionType.vertical:
        final slideAnimation = Tween<Offset>(
          begin: const Offset(0.0, 0.1),
          end: Offset.zero,
        ).chain(CurveTween(curve: Curves.easeOutCubic)).animate(animation);

        return FadeTransition(
          opacity: fadeAnimation,
          child: SlideTransition(
            position: slideAnimation,
            child: child,
          ),
        );
    }
  }
}

enum SharedAxisTransitionType {
  scaled,
  horizontal,
  vertical,
}
